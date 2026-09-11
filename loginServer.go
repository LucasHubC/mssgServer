package main

import (
	"fmt"
	"mssgServer/config"
	"mssgServer/config/db"
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"mssgServer/infra/rpc"
	"mssgServer/net"
	"mssgServer/server/login"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	log.Init()
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	port := config.File.MustValue("login_server", "port", "8003")
	gatewayURL := config.File.MustValue("gateway", "url", "http://127.0.0.1:9000")

	// ✅ 初始化数据库（失败直接退出）
	if err := db.Init(); err != nil {
		log.Fatal("DB init failed:", err)
	}
	defer db.Close()

	// ✅ 初始化 Redis（可选，失败只 warning，不阻塞启动）
	nodeAddr := fmt.Sprintf("ws://%s:%s", host, port)
	if err := db.InitRedis(); err != nil {
		log.Warn("Redis init failed (multi-node features disabled):", err)
	} else {
		defer db.CloseRedis()
		// 告诉 redis 包"我是谁"，PublishBroadcast 会自动填 FromNode
		rediscache.SetMyNodeAddr(nodeAddr)
		// 后台协程：订阅其他节点发来的 Redis 广播
		go net.ConnManager.SubscribeRemoteBroadcast()
	}

	s := net.NewServer(host + ":" + port)
	// 给 Server 设置 NodeAddr，后续每个新连接都会继承这个值
	s.SetNodeAddr(fmt.Sprintf("ws://%s:%s", host, port))
	router := net.NewRouter()
	login.Init(router)
	s.Router(router)

	// ──── 🆕 向 Gateway 注册 + 定期心跳 ────
	httpAddr := fmt.Sprintf("http://%s:%s", host, port) // 这里可以和 ws 共用同一端口
	if err := rpc.Register(gatewayURL, rpc.NodeRegisterReq{
		Type:     "login",
		Addr:     nodeAddr,
		HTTPAddr: httpAddr,
		MaxConns: net.ConnManager.MaxConns(),
	}); err != nil {
		log.Warn("register to gateway failed (is gateway running?):", err)
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
// 收到信号后：广播通知客户端 → 等所有 handler 跑完 → 关 DB → 退出
func gracefulShutdown(done chan struct{}) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// 等信号
	sig := <-sigCh
	log.Info("received shutdown signal", "signal", sig)

	// ① 广播通知所有在线玩家（可选，给客户端一个友好提示）
	log.Info("broadcasting shutdown notice...")
	net.ConnManager.Broadcast("server_msg", map[string]string{
		"type":    "server_shutdown",
		"message": "server is shutting down, please reconnect later",
	})

	// ② 给客户端一点时间收到广播
	time.Sleep(500 * time.Millisecond)

	// ③ 关闭所有连接（触发各自的 readMsgLoop 退出 → Unregister）
	log.Info("closing all connections, total=", net.ConnManager.Total())
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

	// ④ 关数据库和 Redis（由 defer 自动关，如果这里手动关确保顺序也可以）
	// db.Close() 和 db.CloseRedis() 已通过 main() 的 defer 注册

	log.Info("graceful shutdown complete")
	close(done)
	fmt.Println("👋 bye~")
	os.Exit(0)
}
