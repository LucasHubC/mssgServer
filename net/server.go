package net

import (
	"mssgServer/infra/log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	address string
	router  *Router
}

func NewServer(addr string) *Server {
	return &Server{
		address: addr,
	}
}

func (s *Server) Router(r *Router) {
	s.router = r
}

// WebSocket 心跳参数
const (
	heartbeatInterval = 30 * time.Second // 服务器每 30s 发一次 Ping
	heartbeatTimeout  = 90 * time.Second // 90s 没收到 Pong 就断连（3 个心跳周期）
)

func (s *Server) Start() {
	// ⚠️ ReadTimeout/WriteTimeout/IdleTimeout 对 WebSocket 长连接必须设为 0
	// 否则 WebSocket 握手完成后，10s 内没发新请求就会被 HTTP Server 掐掉
	srv := &http.Server{
		Addr:           s.address,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/", s.wsHandler)

	log.Info("server starting", "addr", s.address)
	if err := srv.ListenAndServe(); err != nil {
		log.Error("listen and serve failed", "addr", s.address, "err", err)
	}
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	// ① Upgrade 之前先查连接数，满了返回 503
	if ConnManager.Total() >= ConnManager.MaxConns() {
		http.Error(w, "server full", http.StatusServiceUnavailable)
		return
	}

	// ② 升级 WebSocket
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("websocket upgrade failed", "remote", r.RemoteAddr, "err", err)
		return
	}

	// ③ 注册到全局连接管理器
	wsServer := newWsServer(wsConn)
	wsServer.Router(s.router)
	ConnManager.Register(wsServer)

	// ──── 心跳设置（三要素缺一不可）────
	// ① PongHandler：收到 Pong 后重置 ReadDeadline
	wsConn.SetPongHandler(func(appData string) error {
		wsConn.SetReadDeadline(time.Now().Add(heartbeatTimeout))
		return nil
	})
	// ② PingHandler：兼容客户端主动发 Ping
	wsConn.SetPingHandler(func(appData string) error {
		wsConn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		return wsConn.WriteMessage(websocket.PongMessage, nil)
	})
	// ③ 初始 ReadDeadline + 启动服务器主动发 Ping 的协程
	wsConn.SetReadDeadline(time.Now().Add(heartbeatTimeout))
	go wsHeartbeat(wsConn, wsServer.done)

	defer ConnManager.Unregister(wsServer)

	log.Info("client connected", "addr", wsConn.RemoteAddr().String(), "total", ConnManager.Total())

	// ④ 启动收发循环（非阻塞）
	wsServer.Start()

	// ⑤ 阻塞等连接关闭 —— 等 readMsgLoop/Close 把 done channel 关掉
	// 这样 defer Unregister 会在连接真正断开后才执行
	wsServer.WaitClose()

	log.Info("client disconnected", "addr", wsConn.RemoteAddr().String(), "total", ConnManager.Total())
}

// wsHeartbeat 服务器主动发 Ping 的协程
// 每隔 heartbeatInterval 发一次 Ping，如果写失败说明连接已死，由 readMsgLoop 那边感知到 ReadDeadline 超时后自然断开
func wsHeartbeat(wsConn *websocket.Conn, done <-chan struct{}) {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			wsConn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-done: // ← readMsgLoop/Close 通知了，马上退
			return
		}
	}
}
