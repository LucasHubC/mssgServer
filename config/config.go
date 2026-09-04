package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Unknwon/goconfig"
)

const configFilePath = "conf/conf.ini"

var File *goconfig.ConfigFile

func init() {

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	configPath := currentDir + "/" + configFilePath

	File, err = goconfig.LoadConfigFile(configPath)

	if err != nil {
		log.Fatal("LoadConfigFile failed: " + err.Error())
	}

	//fmt.Println(File)

}

func A() {
	fmt.Println("Hello, World!")
}
