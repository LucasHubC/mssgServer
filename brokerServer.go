package main

import (
	"fmt"
	"mssgServer/config"
	"mssgServer/config/db"
	"mssgServer/gateway"
	"mssgServer/infra/log"
	"mssgServer/infra/metrics"
	rediscache "mssgServer/infra/redis"
	"mssgServer/infra/rpc"
	"net/http"
	_ "net/http/pprof"
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
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("[Broker] listen failed", "err", err)
	}
}
