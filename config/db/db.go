package db

import (
	"database/sql"
	"fmt"
	"mssgServer/config"
	"mssgServer/infra/log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB 全局数据库连接池（业务库 mssg）
var DB *sql.DB

// DBLog 日志库连接池（game_log，独立库，存操作日志等）
// 没有日志库时可以为 nil，DAO 层需要判空
var DBLog *sql.DB

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

	// ✅ 初始化日志库（可选，不配 [mysql_log] 就跳过）
	InitLogDB()
	return nil
}

// InitLogDB 初始化日志库连接（game_log）
// 如果 conf.ini 没有 [mysql_log] section 或连不上，DBLog 保持 nil，不影响主流程
func InitLogDB() {
	// 检查有没有配 [mysql_log]（goconfig 没有 HasSection，用 GetSectionList 检查）
	sections := config.File.GetSectionList()
	hasLogSection := false
	for _, s := range sections {
		if s == "mysql_log" {
			hasLogSection = true
			break
		}
	}
	if !hasLogSection {
		log.Info("[mysql_log] not configured, DBLog is nil")
		return
	}

	host := config.File.MustValue("mysql_log", "host", "127.0.0.1")
	port := config.File.MustValue("mysql_log", "port", "3306")
	user := config.File.MustValue("mysql_log", "user", "root")
	password := config.File.MustValue("mysql_log", "password", "")
	dbname := config.File.MustValue("mysql_log", "dbname", "game_log")
	maxOpen := config.File.MustInt("mysql_log", "max_open", 50)
	maxIdle := config.File.MustInt("mysql_log", "max_idle", 10)
	maxLifetime := config.File.MustInt("mysql_log", "max_lifetime", 300)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname,
	)

	var err error
	DBLog, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Warn("log db sql.Open failed, DBLog will be nil", "err", err)
		return
	}
	DBLog.SetMaxOpenConns(maxOpen)
	DBLog.SetMaxIdleConns(maxIdle)
	DBLog.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

	if err := DBLog.Ping(); err != nil {
		log.Warn("log db Ping failed, DBLog will be nil", "err", err, "host", host, "dbname", dbname)
		DBLog.Close()
		DBLog = nil
		return
	}

	log.Info("log database connected", "host", host, "dbname", dbname, "max_open", maxOpen)
}

// Close 关闭数据库连接池，程序退出前调用
func Close() {
	if DB != nil {
		DB.Close()
		log.Info("database closed")
	}
	if DBLog != nil {
		DBLog.Close()
		log.Info("log database closed")
	}
}
