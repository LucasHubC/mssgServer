// apiLoginServer —— 纯 HTTP API 入口（短连接模式）
// 和 wsLoginServer（loginServer.go）职责完全分离
// 适合：注册、登录、查询、配置读取（不需要推送的接口）

package main

import (
	"context"
	"mssgServer/config"
	"mssgServer/config/db"
	"mssgServer/infra/log"
	"mssgServer/server/inventory"
	"mssgServer/server/login/controller"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Init("apiLogin")

	host := config.File.MustValue("api_server", "host", "127.0.0.1")
	port := config.File.MustValue("api_server", "port", "8006")

	if err := db.Init(); err != nil {
		log.Fatal("DB init failed", "err", err)
	}
	defer db.Close()

	if err := db.InitRedis(); err != nil {
		log.Warn("Redis init failed (cache disabled)", "err", err)
	} else {
		defer db.CloseRedis()
	}

	mux := http.NewServeMux()
	controller.DefaultAccount.RegisterHttpRoutes(mux)
	inventory.InitHttp(mux)
	inventory.StartAutoSave() // HTTP 进程也要启动定时落盘

	srv := &http.Server{
		Addr:         host + ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		log.Info("apiLoginServer shutting down...")
		inventory.Shutdown() // 强制保存 Redis 中脏背包到 MySQL
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
		log.Info("apiLoginServer bye~")
	}()

	log.Info("[apiLoginServer] starting", "addr", host+":"+port)
	log.Info("[apiLoginServer] HTTP routes: /api/login, /api/register, /api/user/{uid}, /api/me, /api/logout")
	log.Info("[apiLoginServer] inventory routes: /api/inventory/add|consume|get")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("apiLoginServer listen failed", "err", err)
	}
}
