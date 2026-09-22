package main

import (
	"fmt"
	"mssgServer/config"
	"mssgServer/config/db"
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"mssgServer/infra/rpc"
	"mssgServer/net"
	"mssgServer/server/chat"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	log.Init("chat")
	host := config.File.MustValue("chat_server", "host", "127.0.0.1")
	port := config.File.MustValue("chat_server", "port", "8004")
	gatewayURL := config.File.MustValue("gateway", "url", "http://127.0.0.1:9000")

	nodeAddr := fmt.Sprintf("ws://%s:%s", host, port)

	// ✅ Redis 是必须的（跨服 Pub/Sub + 群组）
	if err := db.InitRedis(); err != nil {
		log.Fatal("Redis init failed (chat server requires Redis)", "err", err)
	}
	defer db.CloseRedis()

	rediscache.SetMyNodeAddr(nodeAddr)
	go net.ConnManager.SubscribeRemoteBroadcast()
	go net.ConnManager.StartHeartbeatLoop()

	// ✅ 初始化跨服 RPC 模块（Redis Stream + Pub/Sub）
	rpc.Init(nodeAddr)
	// 注册本节点能响应的 RPC 方法
	rpc.Handle("get_online_count", func(p map[string]any) (any, error) {
		return net.ConnManager.Total(), nil
	})
	rpc.Handle("ping", func(p map[string]any) (any, error) {
		return "pong", nil
	})

	s := net.NewServer(host + ":" + port)
	s.SetNodeAddr(nodeAddr)

	router := net.NewRouter()
	chat.Init(router)
	s.Router(router)

	// 向 Gateway 注册，类型 "chat"
	httpAddr := fmt.Sprintf("http://%s:%s", host, port)
	if err := rpc.Register(gatewayURL, rpc.NodeRegisterReq{
		Type:     "chat",
		Addr:     nodeAddr,
		HTTPAddr: httpAddr,
		MaxConns: net.ConnManager.MaxConns(),
	}); err != nil {
		log.Warn("register to gateway failed (is gateway running?)", "err", err)
	}

	// 心跳协程
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			rpc.Heartbeat(gatewayURL, rpc.NodeHeartbeatReq{
				Type:   "chat",
				Addr:   nodeAddr,
				Online: net.ConnManager.Total(),
			})
		}
	}()

	// 优雅关闭
	waitCh := make(chan struct{})
	go gracefulShutdown(waitCh)

	log.Info("[ChatServer] starting on", nodeAddr, "type=chat")
	s.Start()

	<-waitCh
}

func gracefulShutdown(done chan struct{}) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	log.Info("[ChatServer] received shutdown signal", "signal", sig)

	net.ConnManager.Broadcast("chat.system", map[string]string{
		"type":    "server_shutdown",
		"message": "chat server is shutting down",
	})
	time.Sleep(300 * time.Millisecond)

	log.Info("[ChatServer] closing all connections", "total", net.ConnManager.Total())
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

	// ③ 停掉 RPC 模块（不再响应新请求）
	rpc.Close()

	// ④ 等所有正在执行的 handler 跑完
	log.Info("[ChatServer] waiting handlers to finish...")
	net.WaitHandlers()

	log.Info("[ChatServer] graceful shutdown complete")
	close(done)
	fmt.Println("👋 chatServer bye~")
	os.Exit(0)
}
