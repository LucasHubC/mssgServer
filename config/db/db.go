package db

import (
	"database/sql"
	"fmt"
	"mssgServer/config"
	"mssgServer/infra/log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB 全局数据库连接池
var DB *sql.DB

// Init 从配置文件读取数据库参数并初始化连接池
// 在 main() 中调用一次即可
func Init() error {
	host := config.File.MustValue("mysql", "host", "127.0.0.1")
	port := config.File.MustValue("mysql", "port", "3306")
	user := config.File.MustValue("mysql", "user", "root")
	password := config.File.MustValue("mysql", "password", "")
	dbname := config.File.MustValue("mysql", "dbname", "mssg")
	maxOpen := config.File.MustInt("mysql", "max_open", 100)
	maxIdle := config.File.MustInt("mysql", "max_idle", 10)
	maxLifetime := config.File.MustInt("mysql", "max_lifetime", 300)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname,
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open failed: %w", err)
	}

	// 连接池调优
	DB.SetMaxOpenConns(maxOpen)
	DB.SetMaxIdleConns(maxIdle)
	DB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

	// 验证连接
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("db.Ping failed: %w", err)
	}

	log.Info("database connected", "host", host, "dbname", dbname, "max_open", maxOpen)
	return nil
}

// Close 关闭数据库连接池，程序退出前调用
func Close() {
	if DB != nil {
		DB.Close()
		log.Info("database closed")
	}
}
