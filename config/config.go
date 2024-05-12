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
}

var ConfigPtr = &Config{}

func init() {
	ConfigPtr = ReadConfig()
}

func ReadConfig() *Config {
	configData, _ := os.ReadFile("config.json")
	//var config *Config = &Config{}
	err := json.Unmarshal(configData, ConfigPtr)
	if err != nil {
		fmt.Println(err)
		return nil
	} else {
		return ConfigPtr
	}
}

func UpdateConfig(configPtr *Config) {
	data, _ := json.MarshalIndent(configPtr, "", "")
	err := os.WriteFile("config.json", data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
