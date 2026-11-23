package main

import (
	"context"
	"fmt"
	"mssgServer/config"
	"mssgServer/config/db"
	"mssgServer/gateway"
	"mssgServer/infra/log"
	"mssgServer/infra/metrics"
	rediscache "mssgServer/infra/redis"
	"mssgServer/infra/rpc"
	"mssgServer/infra/util"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Broker — 两阶段连接的第二阶段：纯 WS 智能路由代理
// 客户端 → Gateway(/gateway/dispatch) → broker地址 → 客户端连 Broker(/ws)
// Broker → Gateway(/gateway/route) → 拿 login/chat/room 地址 → dialAll
func main() {
	log.Init("broker")

	gatewayURL := config.File.MustValue("broker", "gateway_url", "http://127.0.0.1:9000")
	host := config.File.MustValue("broker", "host", "0.0.0.0")
	port := config.File.MustValue("broker", "port", "9001")
	nodeAddr := fmt.Sprintf("ws://%s:%s/ws", host, port)

	if err := db.InitRedis(); err != nil {
		log.Warn("[Broker] redis init failed", "err", err)
	} else {
		defer db.CloseRedis()
		rediscache.SetMyNodeAddr(fmt.Sprintf("broker:%s:%s", host, port))
	}

	gwClient := gateway.NewHTTPGatewayClient(gatewayURL)

	// Metrics + pprof + health
	metrics.Init()
	metrics.RegisterHTTP()

	// Subscribe Redis kick — admin/kick 需要踢 client WS（做法 A 后 kick backend WS 只会触发重连）
	gateway.StartKickListener()

	// Register 到 Gateway（type=broker）
	nr := rpc.NodeRegisterReq{
		Type:     "broker",
		Addr:     nodeAddr,
		HTTPAddr: "http://" + host + ":" + port,
		MaxConns: 5000,
	}
	if err := rpc.Register(gatewayURL, nr); err != nil {
		log.Fatal("[Broker] ❌ register failed", "err", err)
	}

	// 心跳协程
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("[Broker] heartbeat goroutine PANIC", "err", r)
			}
		}()
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			rpc.Heartbeat(gatewayURL, rpc.NodeHeartbeatReq{
				Type:   nr.Type,
				Addr:   nr.Addr,
				Online: gateway.ActiveProxyConns(),
			})
		}
	}()

	// WS 智能路由代理（复用 gateway.HandleWSProxy）
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		gateway.HandleWSProxy(w, r, "Broker", gwClient)
	})

	addr := host + ":" + port
	log.Info("[Broker] starting on", addr, "→ gateway:", gatewayURL)
	log.Info("[Broker] WS proxy: ws://" + addr + "/ws?token=xxx")

	srv := &http.Server{
		Addr:           addr,
		Handler:        nil,
		MaxHeaderBytes: 1 << 20,           // 1MB header
		ReadTimeout:    10 * time.Second,  // 防慢连接
		WriteTimeout:   10 * time.Second,  // 防慢响应
		IdleTimeout:    120 * time.Second, // Keep-Alive 空闲超时
	}

	// 优雅关闭：监听 Ctrl+C / SIGTERM
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		sig := <-sigCh
		log.Info("[Broker] received shutdown signal", "signal", sig)

		// ① 主动从 Gateway 注销
		log.Info("[Broker] deregistering from gateway...")
		rpc.Deregister(gatewayURL, rpc.NodeHeartbeatReq{
			Type: "broker",
			Addr: nodeAddr,
		})
		time.Sleep(200 * time.Millisecond)

		// ② 停止接受新连接，等最多 5s 让现有代理连接处理完
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Warn("[Broker] graceful shutdown timed out", "err", err)
		}
		log.Info("[Broker] shutdown complete, active proxies=", gateway.ActiveProxyConns())
	}()

	// 自定义 TCP Listener：增大 backlog，抗并发连接风暴
	// util.ListenTCP 直接调 syscall.Listen(fd, 4096)，跨所有 Go 版本生效，附带 SO_REUSEADDR
	ln, err := util.ListenTCP(addr, 4096)
	if err != nil {
		log.Fatal("[Broker] listen failed", "err", err)
	}
	log.Info("[Broker] TCP backlog=4096")

	if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		log.Fatal("[Broker] serve failed", "err", err)
	}
	log.Info("[Broker] bye~")
}
