package net

import (
	"encoding/json"
	"mssgServer/infra/log"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

// WSServer 服务器端
type WSServer struct {
	wsConn *websocket.Conn

	router *Router

	outChan chan *WsMsgResp

	Seq int64

	property map[string]interface{}

	propertyLock sync.RWMutex

	// limiter 每连接令牌桶限流器
	// 每秒最多 30 条消息，突发最多 50 条
	limiter *rate.Limiter

	// done 通知 wsHandler 连接已关闭
	done chan struct{}

	// once 保证 Close 只执行一次
	once sync.Once
}

func newWsServer(wsConn *websocket.Conn) *WSServer {
	return &WSServer{
		wsConn:   wsConn,
		outChan:  make(chan *WsMsgResp, 256),
		property: make(map[string]interface{}),
		Seq:      0,
		limiter:  rate.NewLimiter(30, 50),
		done:     make(chan struct{}),
	}
}

func (s *WSServer) Router(router *Router) {
	s.router = router

}

func (s *WSServer) SetProperty(key string, value interface{}) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.property[key] = value
}

func (s *WSServer) GetProperty(key string) (interface{}, error) {
	s.propertyLock.RLock()
	defer s.propertyLock.RUnlock()
	if value, ok := s.property[key]; ok {
		return value, nil
	}
	return nil, nil
}

func (s *WSServer) RemoveProperty(key string) {
	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	delete(s.property, key)
}
func (s *WSServer) Addr() string {
	return s.wsConn.RemoteAddr().String()
}

func (s *WSServer) Push(name string, data interface{}) {

	rsp := &WsMsgResp{
		Body: &RespBody{
			Seq:  s.Seq,
			Name: name,
			Code: 0,
			Msg:  data,
		},
	}
	s.outChan <- rsp
}

// 通道一旦建立 那么收发消息都要持续去监听通道
func (s *WSServer) Start() {
	go s.writeMsgLoop()
	go s.readMsgLoop()
}

func (s *WSServer) writeMsgLoop() {
	for msg := range s.outChan {
		data, err := json.Marshal(msg.Body)
		if err != nil {
			log.Error("marshal response failed", "err", err)
			continue
		}
		if err := s.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Error("write message failed", "addr", s.Addr(), "err", err)
			return
		}
	}
}

func (s *WSServer) readMsgLoop() {
	//先读消息 再处理消息

	defer func() {
		if err := recover(); err != nil {
			log.Error("readMsgLoop panic", "addr", s.Addr(), "err", err)
			s.Close()
		}
	}()

	for {
		_, messageData, err := s.wsConn.ReadMessage() // 阻塞等待客户端发来消息

		if err != nil {
			log.Warn("read message failed", "addr", s.Addr(), "err", err)
			break
		}

		log.Debug("recv raw", "addr", s.Addr(), "data", string(messageData))

		body := &ReqBody{}
		err = json.Unmarshal(messageData, body)
		if err != nil {
			log.Warn("unmarshal message failed", "addr", s.Addr(), "err", err, "raw", string(messageData))
			break
		}

		// ✅ 限流检查：Allow 不阻塞，超限直接断开
		if !s.limiter.Allow() {
			log.Warn("rate limit exceeded, disconnecting", "addr", s.Addr())
			break
		}

		// 路由
		req := &WsMsgReq{Body: body, Conn: s}
		rsp := &WsMsgResp{Body: &RespBody{Name: body.Name, Seq: req.Body.Seq}}
		s.router.Run(req, rsp)
		s.outChan <- rsp

	}

	s.Close()
}

// Close 关闭连接，幂等（sync.Once 保证只执行一次）
// 关闭顺序：先 close outChan（让 writeMsgLoop 退出），再关 wsConn
func (s *WSServer) Close() {
	s.once.Do(func() {
		close(s.outChan)
		s.wsConn.Close()
		close(s.done)
	})
}

// WaitClose 阻塞等连接关闭（wsHandler 里用）
func (s *WSServer) WaitClose() {
	<-s.done
}
