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

func (s *Server) Start() {
	srv := &http.Server{
		Addr:           s.address,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
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
	defer ConnManager.Unregister(wsServer)

	log.Info("client connected", "addr", wsConn.RemoteAddr().String(), "total", ConnManager.Total())

	// ④ 启动收发循环（非阻塞）
	wsServer.Start()

	// ⑤ 阻塞等连接关闭 —— 等 readMsgLoop/Close 把 done channel 关掉
	// 这样 defer Unregister 会在连接真正断开后才执行
	wsServer.WaitClose()

	log.Info("client disconnected", "addr", wsConn.RemoteAddr().String(), "total", ConnManager.Total())
}
