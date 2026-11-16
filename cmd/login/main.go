package main

import (
	"fmt"
	"mssgServer/config"
	"mssgServer/config/db"
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"mssgServer/infra/rpc"
	"mssgServer/net"
	"mssgServer/server/inventory"
	"mssgServer/server/login"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// 包级变量 —— main() 和 gracefulShutdown() 都要访问
var (
	gatewayURL string
	nodeAddr   string
	s          *net.Server
)

func main() {

	log.Init("login")
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	port := config.File.MustValue("login_server", "port", "8003")
	gatewayURL = config.File.MustValue("gateway", "url", "http://127.0.0.1:9000")

	// ✅ 初始化数据库（失败直接退出）
	if err := db.Init(); err != nil {
		log.Fatal("DB init failed", "err", err)
	}
	defer db.Close()

	// ✅ 初始化 Redis（可选，失败只 warning，不阻塞启动）
	nodeAddr = fmt.Sprintf("ws://%s:%s", host, port)
	if err := db.InitRedis(); err != nil {
		log.Warn("Redis init failed (multi-node features disabled)", "err", err)
	} else {
		defer db.CloseRedis()
		// 告诉 redis 包"我是谁"，PublishBroadcast 会自动填 FromNode
		rediscache.SetMyNodeAddr(nodeAddr)
		// 后台协程：订阅其他节点发来的 Redis 广播
		go net.ConnManager.SubscribeRemoteBroadcast()
		// 后台协程：定时续期本机所有在线用户的 Redis 在线状态（防止 30s 过期被误判下线）
		go net.ConnManager.StartHeartbeatLoop()

		// ✅ 初始化跨服 RPC 模块（Redis Stream + Pub/Sub）
		rpc.Init(nodeAddr)
		// 注册本节点能响应的 RPC 方法 —— 其他节点可以 Call() 过来
		rpc.Handle("get_online_count", func(p map[string]any) (any, error) {
			return net.ConnManager.Total(), nil
		})
		rpc.Handle("ping", func(p map[string]any) (any, error) {
			return "pong", nil
		})

		// ────────────────────────────────────
		// 🆕 get_inventory —— 跨服查背包
		// ────────────────────────────────────
		// 其他节点（chatServer / gameServer）可以 Call() 过来查指定玩家的道具
		// 数据走 Redis 热数据层，不碰 DB，高性能
		rpc.Handle("get_inventory", func(p map[string]any) (any, error) {
			uid, _ := p["uid"].(string)
			if uid == "" {
				return nil, fmt.Errorf("uid is required")
			}
			items, err := rediscache.GetAllItems(uid)
			if err != nil {
				return nil, fmt.Errorf("redis read failed: %w", err)
			}
			return map[string]any{
				"uid":   uid,
				"items": items,
			}, nil
		})

		// ────────────────────────────────────
		// 🆕 [演示] 跨服 RPC 自测
		// ────────────────────────────────────
		// 启动 6 秒后自动 ping chatServer，验证跨服 RPC 链路通不通
		go func() {
			time.Sleep(6 * time.Second)

			// 从配置文件读 chatServer 地址
			chatAddr := fmt.Sprintf("ws://%s:%s",
				config.File.MustValue("chat_server", "host", "127.0.0.1"),
				config.File.MustValue("chat_server", "port", "8004"))

			// 测试 1：loginServer → chatServer ping
			result, err := rpc.Call(chatAddr, "ping", nil, 3*time.Second)
			if err != nil {
				log.Warn("[RPC Demo] loginServer → chatServer ping failed",
					"err", err, "💡 可能 chatServer 没启动")
			} else {
				log.Info("[RPC Demo] ✅ loginServer → chatServer ping 通了！",
					"result", result)
			}

			// 测试 2：调自己的 get_inventory（演示 handler 工作正常）
			myResult, myErr := rpc.Call(nodeAddr, "get_inventory",
				map[string]any{"uid": "_any_demo_uid_"}, 3*time.Second)
			log.Info("[RPC Demo] 调自己 get_inventory",
				"err", myErr, "result", myResult)
		}()

	}

	s = net.NewServer(host + ":" + port)
	// 给 Server 设置 NodeAddr，后续每个新连接都会继承这个值
	s.SetNodeAddr(fmt.Sprintf("ws://%s:%s", host, port))
	router := net.NewRouter()
	// 全局中间件：Recover 兜底 panic、Log 打印请求耗时
	router.Use(net.RecoverMiddleware, net.LogMiddleware)
	login.Init(router)
	inventory.Init(router) // ← 道具模块：WebSocket 路由 + 定时落盘 goroutine
	s.Router(router)

	// ──── 🆕 向 Gateway 注册 + 定期心跳 ────
	httpAddr := fmt.Sprintf("http://%s:%s", host, port) // 这里可以和 ws 共用同一端口
	if err := rpc.Register(gatewayURL, rpc.NodeRegisterReq{
		Type:     "login",
		Addr:     nodeAddr,
		HTTPAddr: httpAddr,
		MaxConns: net.ConnManager.MaxConns(),
	}); err != nil {
		log.Warn("register to gateway failed (is gateway running?)", "err", err)
	}
	// 心跳协程：每 5s 上报一次当前在线人数
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			rpc.Heartbeat(gatewayURL, rpc.NodeHeartbeatReq{
				Type:   "login",
				Addr:   nodeAddr,
				Online: net.ConnManager.Total(),
			})
		}
	}()

	// ──── 优雅关闭：监听 Ctrl+C / SIGTERM ────
	// 必须放在 s.Start() 之前，因为 Start() 内部 ListenAndServe 是阻塞的
	waitCh := make(chan struct{}) // 关闭完成信号
	go gracefulShutdown(waitCh)

	s.Start() // 阻塞在这里，直到进程退出

	<-waitCh // 等 gracefulShutdown 收尾完成后再彻底退出

}

// gracefulShutdown 监听 OS 信号，做优雅关闭
// 收到信号后：先停 accept → 广播通知 → 关所有 WS → 等 handler → 退出
func gracefulShutdown(done chan struct{}) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// 等信号
	sig := <-sigCh
	log.Info("received shutdown signal", "signal", sig)

	// ① 先停 HTTP accept（net.Server 内部 ListenAndServe 是 goroutine 非阻塞的）
	//    这步之后不会再接受新的 WebSocket 连接
	log.Info("stopping HTTP server...")
	s.Shutdown()

	// ② 主动从 Gateway 注销 —— 让 Gateway 立即摘除本节点
	log.Info("deregistering from gateway...")
	rpc.Deregister(gatewayURL, rpc.NodeHeartbeatReq{
		Type: "login",
		Addr: nodeAddr,
	})
	// 给 Gateway 一点时间处理
	time.Sleep(200 * time.Millisecond)

	// ③ 广播通知所有在线玩家（给客户端一个友好提示）
	log.Info("broadcasting shutdown notice...")
	net.ConnManager.Broadcast("server_msg", map[string]string{
		"type":    "server_shutdown",
		"message": "server is shutting down, please reconnect later",
	})

	// ④ 给客户端一点时间收到广播
	time.Sleep(500 * time.Millisecond)

	// ⑤ 关闭所有连接（触发各自的 readMsgLoop 退出 → Unregister）
	log.Info("closing all connections", "total", net.ConnManager.Total())
	all := net.ConnManager.All()
	var wg sync.WaitGroup
	for _, ws := range all {
		wg.Add(1)
		go func(ws *net.WSServer) {
			defer wg.Done()
			ws.Close()
		}(ws)
	}
	wg.Wait()

	// ⑤ 等所有正在执行的 handler 跑完（防止 handler 正在写 DB 时被强制断开）
	log.Info("waiting handlers to finish...")
	net.WaitHandlers()

	// ⑥ 停掉 RPC 模块（不再响应新请求）
	rpc.Close()

	// ⑦ 强制保存所有脏背包到 MySQL（防止最后 30s 内的修改丢失）
	log.Info("saving dirty inventory...")
	inventory.Shutdown()

	log.Info("graceful shutdown complete")
	close(done)
	log.Info("loginServer bye~")
	os.Exit(0)
}
