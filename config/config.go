package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	TargetDir string
	Width     int
	Height    int
	Title     string
	WhiteList []string
}

var ConfigPtr = &Config{}

// init 导入包时初始化 ConfigPtr
func init() {
	fmt.Println("init config")
	ConfigPtr.readConfig()
}

// 重新读取 config.json 文件
func (configPtr *Config) readConfig() {
	configData, _ := os.ReadFile("config/config.json")
	//var config *Config = &Config{}
	err := json.Unmarshal(configData, configPtr)
	if err != nil {
		fmt.Println("Failed init config,", err)
	}
}

// 将 ConfigPtr 写入 config.json
func (configPtr *Config) writeConfig() {
	data, _ := json.MarshalIndent(configPtr, "", "")
	err := os.WriteFile("config/config.json", data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
