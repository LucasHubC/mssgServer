package chat

import (
	"mssgServer/net"
	"mssgServer/server/chat/controller"
)

func Init(router *net.Router) {
	initRouter(router)
}

func initRouter(router *net.Router) {
	controller.DefaultChat.Router(router)
}
