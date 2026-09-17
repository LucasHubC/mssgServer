package net

import (
	"mssgServer/infra/log"
	rediscache "mssgServer/infra/redis"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	address  string
	router   *Router
	nodeAddr string // 当前节点的 ws 地址（如 "ws://127.0.0.1:8003"），跨节点 Redis 同步用
}

func NewServer(addr string) *Server {
	return &Server{
		address: addr,
	}
}

func (s *Server) Router(r *Router) {
	s.router = r
}

// SetNodeAddr 设置当前节点的 ws 地址（loginServer 启动时调一次）
func (s *Server) SetNodeAddr(addr string) {
	s.nodeAddr = addr
}

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
		WriteJSONWithStatus(w, http.StatusServiceUnavailable,
			map[string]any{"code": -1, "errMsg": "server full"})
		return
	}

	// ② 升级 WebSocket
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("websocket upgrade failed", "remote", r.RemoteAddr, "err", err)
		return
	}

	// ③ 创建 WSServer
	wsServer := newWsServer(wsConn)
	wsServer.Router(s.router)
	wsServer.NodeAddr = s.nodeAddr

	// ④ Query Token 自动鉴权（跨 Server 统一支持）
	// 客户端连 chatServer 时可以 ws://host:8004?token=xxx，自动从 Redis session 绑定 uid
	if token := r.URL.Query().Get("token"); token != "" {
		if sess, ok := rediscache.GetSession(token); ok {
			wsServer.SetProperty("uid", sess.UId)
			wsServer.SetProperty("username", sess.Username)
			log.Info("ws auto auth via token", "uid", sess.UId, "username", sess.Username)
			// 绑定到 ConnManager（跨服 PushToUid / BindUid 需要）
			ConnManager.BindUid(sess.UId, wsServer)
		} else {
			log.Warn("ws auth failed: invalid token", "remote", r.RemoteAddr)
		}
	}

	// ⑤ 注册到全局连接管理器
	ConnManager.Register(wsServer)

	defer ConnManager.Unregister(wsServer)

	log.Info("client connected", "addr", wsConn.RemoteAddr().String(), "total", ConnManager.Total())

	// ⑥ 启动收发循环（非阻塞）
	wsServer.Start()

	// ⑦ 给客户端发鉴权结果（Start 之后才能 Push，outChan 已初始化）
	if uid := GetUid(wsServer); uid != "" {
		wsServer.Push("auth.ok", map[string]any{
			"uid":      uid,
			"username": GetUsername(wsServer),
		})
	} else if r.URL.Query().Get("token") != "" {
		wsServer.Push("auth.fail", map[string]any{
			"errMsg": "token 无效或已过期",
		})
	}

	// ⑤ 阻塞等连接关闭 —— 等 readMsgLoop/Close 把 done channel 关掉
	// 这样 defer Unregister 会在连接真正断开后才执行
	wsServer.WaitClose()

	log.Info("client disconnected", "addr", wsConn.RemoteAddr().String(), "total", ConnManager.Total())
}
