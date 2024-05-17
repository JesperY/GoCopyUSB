package config

import (
	"encoding/json"
	"fmt"
	"gioui.org/app"
	"github.com/JesperY/GoCopyUSB/logger"
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
	Win               *app.Window
	AutoStartUp       bool
}

var ConfigPtr = &Config{}

// init 导入包时初始化 ConfigPtr
func init() {
	fmt.Println("init config")
	ConfigPtr.ReadConfig()
}

func (configPtr *Config) setDefault() {
	configPtr.TargetDir = "."
	configPtr.Width = 600
	configPtr.Height = 400
	configPtr.Title = "USBCopier"
	configPtr.DelayMinutes = 0
	configPtr.WhiteListDir = []string{}
	configPtr.WhiteListFilename = []string{}
	configPtr.WhiteListSuffix = []string{}
	configPtr.AutoStartUp = false
	configPtr.WriteConfig()
}

// ReadConfig 重新读取 config.json 文件
func (configPtr *Config) ReadConfig() {
	configData, _ := os.ReadFile("config/config.json")
	//var config *Config = &Config{}
	err := json.Unmarshal(configData, configPtr)
	if err != nil {
		fmt.Println("Failed init config,", err)
		logger.SugarLogger.Errorf("Failed init config, %v, using default setting.", err)
		// 读取配置文件失败，使用采用默认值，并创建配置文件 config.json 将默认值写入
		ConfigPtr.setDefault()
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
		//fmt.Println(err)
		// 无法写入配置文件
		logger.SugarLogger.Errorf("Failed writing config.json, %v, using default setting.", err)
	}
}
func (configPtr *Config) WriteConfigWithPath(path string) {
	data, _ := json.MarshalIndent(configPtr, "", "")
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
