package room

import (
	"mssgServer/net"
	"mssgServer/server/room/controller"
)

// Init 注册房间 WebSocket 路由
//
//	room.create / room.join / room.leave / room.kick / room.delete / room.send
func Init(r *net.Router) {
	controller.DefaultRoom.Router(r)
}
