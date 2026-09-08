package main

import (
	"mssgServer/config"
	"mssgServer/gateway"
	"mssgServer/infra/log"
	"net/http"
)

func main() {
	log.Init()

	host := config.File.MustValue("gateway", "host", "0.0.0.0")
	port := config.File.MustValue("gateway", "port", "9000")
	policy := config.File.MustValue("gateway", "policy", "round_robin")

	gtw := gateway.NewGateway(gateway.RoutePolicy(policy))

	// 注册后端节点的 HTTP API
	http.HandleFunc("/gateway/register", gtw.HandleRegister)
	http.HandleFunc("/gateway/heartbeat", gtw.HandleHeartbeat)
	http.HandleFunc("/gateway/route", gtw.HandleRoute)   // 客户端用这个
	http.HandleFunc("/gateway/status", gtw.HandleStatus) // 调试

	// 启动清理协程
	go gtw.CleanupLoop()

	addr := host + ":" + port
	log.Info("[Gateway] starting on", addr, "policy=", policy)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("[Gateway] listen failed:", err)
	}
}
