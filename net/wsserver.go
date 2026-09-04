package net

import (
	"encoding/json"
	"mssgServer/infra/log"
	"sync"
	"github.com/gorilla/websocket"
)

// WSServer 服务器端
type WSServer struct {
	wsConn *websocket.Conn

	router *Router

	outChan chan *WsMsgResp

	Seq int64

	property map[string]interface{}

	propertyLock sync.RWMutex
}

func newWsServer(wsConn *websocket.Conn) *WSServer {
	return &WSServer{
		wsConn:   wsConn,
		outChan:  make(chan *WsMsgResp),
		property: make(map[string]interface{}),
		Seq:      0,
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
		//路由
		req := &WsMsgReq{Body: body, Conn: s}
		rsp := &WsMsgResp{Body: &RespBody{Name: body.Name, Seq: req.Body.Seq}}
		s.router.Run(req, rsp)
		s.outChan <- rsp

	}

	s.Close()
}

func (s *WSServer) Close() {
	s.wsConn.Close()
}
