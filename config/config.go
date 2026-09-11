package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Unknwon/goconfig"
)

const configFilePath = "conf/conf.ini"

var File *goconfig.ConfigFile

func init() {
	// 从当前目录开始，逐级向上找 conf/conf.ini
	// 这样不管从哪里启动（go run、VS Code 测试、go test ./xxx）都能找到配置
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	configPath := ""
	for {
		candidate := filepath.Join(dir, configFilePath)
		if _, err := os.Stat(candidate); err == nil {
			configPath = candidate
			break
		}
		// 向上一级
		parent := filepath.Dir(dir)
		if parent == dir { // 已经到根目录还没找到
			break
		}
		dir = parent
	}

	if configPath == "" {
		log.Fatal("LoadConfigFile failed: conf/conf.ini not found in any parent directory")
	}

	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		log.Fatal("LoadConfigFile failed: " + err.Error())
	}
}
