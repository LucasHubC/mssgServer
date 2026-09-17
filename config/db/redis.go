package db

import (
	"context"
	"fmt"
	"mssgServer/config"
	"mssgServer/infra/log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RDB 全局 Redis 客户端
var RDB *redis.Client

// InitRedis 从配置文件读取 Redis 参数并初始化连接池
// Redis 是可选的：如果连不上只打 warning，不阻塞启动
// 业务层调用前应该先判断 RDB != nil
func InitRedis() error {
	host := config.File.MustValue("redis", "host", "127.0.0.1")
	port := config.File.MustValue("redis", "port", "6379")
	password := config.File.MustValue("redis", "password", "")
	dbIdx := config.File.MustInt("redis", "db", 0)

	RDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Password:     password,
		DB:           dbIdx,
		PoolSize:     100,
		MinIdleConns: 10,
	})

	// 验证连接
	pingCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := RDB.Ping(pingCtx).Err(); err != nil {
		return fmt.Errorf("redis.Ping failed: %w", err)
	}

	log.Info("redis connected", "addr", fmt.Sprintf("%s:%s", host, port))
	return nil
}

// CloseRedis 关闭 Redis 客户端
func CloseRedis() {
	if RDB != nil {
		RDB.Close()
		log.Info("redis closed")
	}
}
