package main

import (
	"fmt"
	"mssgServer/config"
	"mssgServer/config/db"
	"mssgServer/gateway"
	"mssgServer/infra/log"
	"mssgServer/infra/metrics"
	rediscache "mssgServer/infra/redis"
	"net/http"
	_ "net/http/pprof"
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
	}

	gtw := gateway.NewGateway(gateway.RoutePolicy(policy))

	// ──── Metrics + pprof + health ────
	metrics.Init()
	metrics.RegisterHTTP()

	// ──── 后端节点用的 HTTP API ────
	http.HandleFunc("/gateway/register", gtw.HandleRegister)
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

	addr := host + ":" + port
	log.Info("[Gateway] starting on", addr, "policy=", policy)
	log.Info("[Gateway] WS proxy:  ws://" + addr + "/ws?token=xxx")
	log.Info("[Gateway] admin APIs: /admin/broadcast, /admin/push, /admin/kick")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("[Gateway] listen failed", "err", err)
	}

}
