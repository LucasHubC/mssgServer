package login

import (
	"mssgServer/net"
	"mssgServer/server/login/controller"
)

func Init(router *net.Router) {
	initRouter(router)
}

func initRouter(router *net.Router) {
	controller.DefaultAccount.Router(router)
}
