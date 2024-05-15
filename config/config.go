package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	TargetDir         string
	Width             int
	Height            int
	Title             string
	DelayMinutes      int
	WhiteListDir      []string
	WhiteListFilename []string
	WhiteListSuffix   []string
}

var ConfigPtr = &Config{}

// init 导入包时初始化 ConfigPtr
func init() {
	fmt.Println("init config")
	ConfigPtr.ReadConfig()
}

// ReadConfig 重新读取 config.json 文件
func (configPtr *Config) ReadConfig() {
	configData, _ := os.ReadFile("config/config.json")
	//var config *Config = &Config{}
	err := json.Unmarshal(configData, configPtr)
	if err != nil {
		fmt.Println("Failed init config,", err)
	}
}

func (configPtr *Config) ReadConfigWithPath(path string) {
	configData, _ := os.ReadFile(path)
	//var config *Config = &Config{}
	err := json.Unmarshal(configData, configPtr)
	if err != nil {
		fmt.Println("Failed init config,", err)
	}
}

// WriteConfig 将 ConfigPtr 写入 config.json
func (configPtr *Config) WriteConfig() {
	data, _ := json.MarshalIndent(configPtr, "", "")
	err := os.WriteFile("config/config.json", data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
func (configPtr *Config) WriteConfigWithPath(path string) {
	data, _ := json.MarshalIndent(configPtr, "", "")
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
