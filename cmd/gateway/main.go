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
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Init("gateway")

	host := config.File.MustValue("gateway", "host", "0.0.0.0")
	port := config.File.MustValue("gateway", "port", "9000")
	policy := config.File.MustValue("gateway", "policy", "round_robin")

	// Gateway 不需要 MySQL，但需要 Redis（用于广播消息转发）
	if err := db.InitRedis(); err != nil {
		log.Warn("[Gateway] redis init failed (admin broadcast disabled)", "err", err)
	} else {
		defer db.CloseRedis()
		// Gateway 自己也会发广播，设个特殊 nodeAddr 用来跳过自己
		rediscache.SetMyNodeAddr(fmt.Sprintf("gateway:%s:%s", host, port))
		// admin API 需要 PublishBroadcast → 必须启动异步广播后台协程
		rediscache.StartAsyncBroadcast()
	}

	gtw := gateway.NewGateway(gateway.RoutePolicy(policy))

	// ──── Metrics + pprof + health ────
	metrics.Init()
	metrics.RegisterHTTP()

	// ──── 后端节点用的 HTTP API ────
	http.HandleFunc("/gateway/register", gtw.HandleRegister)
	http.HandleFunc("/gateway/deregister", gtw.HandleDeregister)
	http.HandleFunc("/gateway/heartbeat", gtw.HandleHeartbeat)
	http.HandleFunc("/gateway/route", gtw.HandleRoute)
	http.HandleFunc("/gateway/status", gtw.HandleStatus)

	// ──── 两阶段连接的第一阶段：客户端先调 dispatch 拿 broker 地址 ────
	http.HandleFunc("/gateway/dispatch", gtw.HandleDispatch)

	// ──── 管理后台用的 HTTP API（跨服广播/推送/踢人）────
	http.HandleFunc("/admin/broadcast", gtw.HandleAdminBroadcast)
	http.HandleFunc("/admin/push", gtw.HandleAdminPush)
	http.HandleFunc("/admin/kick", gtw.HandleAdminKick)

	// ──── WS 智能路由代理（Gateway 进程内：gtw 就是 NodeProvider）────
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		gateway.HandleWSProxy(w, r, "Gateway", gtw)
	})

	// 启动清理协程
	go gtw.CleanupLoop()

	// ⚡ 超时配置与 brokerServer / net.Server 保持一致
	addr := host + ":" + port
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
		log.Info("[Gateway] received shutdown signal", "signal", sig)

		// 停止接受新连接，等最多 5s 让 WS proxy 连接处理完
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Warn("[Gateway] graceful shutdown timed out", "err", err)
		} else {
			log.Info("[Gateway] graceful shutdown complete")
		}
	}()

	log.Info("[Gateway] starting on", addr, "policy=", policy)
	log.Info("[Gateway] WS proxy:  ws://" + addr + "/ws?token=xxx")
	log.Info("[Gateway] admin APIs: /admin/broadcast, /admin/push, /admin/kick")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("[Gateway] listen failed", "err", err)
	}

	fmt.Println("👋 gatewayServer bye~")
}
