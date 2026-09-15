// apiLoginServer —— 纯 HTTP API 入口（短连接模式）
// 和 wsLoginServer（loginServer.go）职责完全分离
// 适合：注册、登录、查询、配置读取（不需要推送的接口）

package main

import (
	"context"
	"fmt"
	"mssgServer/config"
	"mssgServer/config/db"
	"mssgServer/infra/log"
	"mssgServer/server/login/controller"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Init()

	host := config.File.MustValue("api_server", "host", "127.0.0.1")
	port := config.File.MustValue("api_server", "port", "8005")

	// ✅ HTTP API 不需要 Redis 广播，但注册/登录走 DAO 需要 MySQL
	if err := db.Init(); err != nil {
		log.Fatal("DB init failed:", err)
	}
	defer db.Close()

	// Redis 可选（DAO 层有 Redis 缓存更快）
	if err := db.InitRedis(); err != nil {
		log.Warn("Redis init failed (cache disabled):", err)
	} else {
		defer db.CloseRedis()
	}

	// 组装路由（纯标准库，零抽象层）
	mux := http.NewServeMux()
	controller.DefaultAccount.RegisterHttpRoutes(mux)

	// HTTP Server —— 参数和 WebSocket Server 完全不同
	srv := &http.Server{
		Addr:         host + ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second, // 短连接：读请求超时严格
		WriteTimeout: 10 * time.Second, // 短连接：写响应超时严格
		IdleTimeout:  30 * time.Second, // 空闲连接回收
	}

	// 优雅关闭（HTTP 的优雅关闭很简单——停止接新请求，等正在处理的完成）
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		log.Info("apiLoginServer shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
		fmt.Println("👋 apiLoginServer bye~")
	}()

	log.Info("[apiLoginServer] starting on", host+":"+port)
	log.Info("[apiLoginServer] routes: POST /api/login, POST /api/register, GET /api/user/{uid}")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("apiLoginServer listen failed:", err)
	}
}
