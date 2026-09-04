package net

import (
	"log"
	"net/http"

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
	log.Println("Server Start: " + s.address)

	http.HandleFunc("/", s.wsHandler)
	err := http.ListenAndServe(s.address, nil)

	if err != nil {
		log.Fatal("ListenAndServe failed: " + err.Error())
	}

}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed: " + err.Error())
		return
	}
	//defer wsConn.Close()

	wsServer := newWsServer(wsConn)
	wsServer.Start()
	wsServer.Router(s.router)

	//websocket 连接通道建立以后 不管是客户端还是服务器 都可以收发消息

	//发消息的是时候 把消息当作路由 来去处理 消息是有格式的 先定义消息的格式 再根据消息的格式去处理消息

	// werr := wsConn.WriteMessage(websocket.TextMessage, []byte("hello world"))
	// if werr != nil {
	// 	log.Println("WriteMessage failed: " + werr.Error())
	// }

	// for {
	// 	messageType, message, err := wsConn.ReadMessage()  // 阻塞等待客户端发来消息
	// 	if err != nil {
	// 		log.Println("ReadMessage failed:", err)
	// 		break  // 出错就跳出循环
	// 	}

	// 	log.Printf("收到消息: %s\n", string(message))

	// 	// 回给客户端
	// 	err = wsConn.WriteMessage(messageType, message)
	// 	if err != nil {
	// 		log.Println("WriteMessage failed:", err)
	// 		break
	// 	}
	// }

}
