package main

import (
	"fmt"
	"mssgServer/config"
	"mssgServer/config/db"
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"mssgServer/infra/rpc"
	"mssgServer/net"
	"mssgServer/server/room"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	log.Init("room")
	host := config.File.MustValue("room_server", "host", "127.0.0.1")
	port := config.File.MustValue("room_server", "port", "8005")
	gatewayURL := config.File.MustValue("gateway", "url", "http://127.0.0.1:9000")

	// roomServer 不需要 MySQL，只需要 Redis
	nodeAddr := fmt.Sprintf("ws://%s:%s", host, port)
	if err := db.InitRedis(); err != nil {
		log.Warn("Redis init failed (multi-node features disabled)", "err", err)
	} else {
		defer db.CloseRedis()
		rediscache.SetMyNodeAddr(nodeAddr)
		go net.ConnManager.SubscribeRemoteBroadcast()
		go net.ConnManager.StartHeartbeatLoop()

		rpc.Init(nodeAddr)
		rpc.Handle("ping", func(p map[string]any) (any, error) {
			return "pong", nil
		})
		rpc.Handle("get_online_count", func(p map[string]any) (any, error) {
			return net.ConnManager.Total(), nil
		})
	}

	s := net.NewServer(host + ":" + port)
	s.SetNodeAddr(nodeAddr)
	router := net.NewRouter()
	room.Init(router)
	s.Router(router)

	// 向 Gateway 注册
	httpAddr := fmt.Sprintf("http://%s:%s", host, port)
	if err := rpc.Register(gatewayURL, rpc.NodeRegisterReq{
		Type:     "room",
		Addr:     nodeAddr,
		HTTPAddr: httpAddr,
		MaxConns: net.ConnManager.MaxConns(),
	}); err != nil {
		log.Warn("register to gateway failed (is gateway running?)", "err", err)
	}
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			rpc.Heartbeat(gatewayURL, rpc.NodeHeartbeatReq{
				Type:   "room",
				Addr:   nodeAddr,
				Online: net.ConnManager.Total(),
			})
		}
	}()

	// 优雅关闭
	waitCh := make(chan struct{})
	go gracefulShutdown(waitCh)

	s.Start()
	<-waitCh
}

func gracefulShutdown(done chan struct{}) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	log.Info("received shutdown signal", "signal", sig)

	net.ConnManager.Broadcast("server_msg", map[string]string{
		"type":    "server_shutdown",
		"message": "roomServer shutting down",
	})
	time.Sleep(500 * time.Millisecond)

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

	// ③ 等所有正在执行的 handler 跑完
	log.Info("waiting handlers to finish...")
	net.WaitHandlers()

	// ④ 停掉 RPC 模块
	rpc.Close()

	// ⑤ 强制退出（因为 s.Start() 内部 ListenAndServe 是阻塞的，main 走不到 <-waitCh）
	log.Info("roomServer graceful shutdown complete")
	fmt.Println("👋 roomServer bye~")
	os.Exit(0)
}
