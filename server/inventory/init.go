package inventory

import (
	"mssgServer/net"
	"mssgServer/server/inventory/controller"
	"net/http"
)

// Init 注册 WebSocket 路由 + 启动定时落盘（loginServer 用）
func Init(router *net.Router) {
	controller.DefaultItem.Router(router)
	StartAutoSave()
}

// InitHttp 仅注册 HTTP REST 路由（apiLoginServer 用）
// 注意：HTTP 进程需要自己调 StartAutoSave()
func InitHttp(mux *http.ServeMux) {
	controller.DefaultItem.RegisterHttpRoutes(mux)
}

// StartAutoSave 启动定时落盘 goroutine（WS 和 HTTP 进程都要调）
func StartAutoSave() {
	controller.DefaultItem.StartAutoSave()
}

// Shutdown 优雅关闭时强制保存所有脏 uid
func Shutdown() {
	controller.DefaultItem.SaveAllDirty()
}
