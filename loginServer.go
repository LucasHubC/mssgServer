package main

import (
	"mssgServer/config"
	"mssgServer/config/db"
	"mssgServer/infra/log"
	"mssgServer/net"
	"mssgServer/server/login"
)

func main() {

	log.Init()
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	port := config.File.MustValue("login_server", "port", "8003")

	// ✅ 初始化数据库（失败直接退出）
	if err := db.Init(); err != nil {
		log.Fatal("DB init failed:", err)
	}
	defer db.Close()

	s := net.NewServer(host + ":" + port)
	router := net.NewRouter()
	login.Init(router)
	s.Router(router)
	s.Start()

}
