package main

import (
	"gioui.org/app"
	"github.com/JesperY/GoCopyUSB/backend"
	"github.com/JesperY/GoCopyUSB/frontend/home"
	"github.com/JesperY/GoCopyUSB/frontend/tray"
	"github.com/JesperY/GoCopyUSB/logger"
)

func main() {
	backend.InitCheck()
	defer logger.SugarLogger.Sync()
	go tray.SysTrayRun()
	go backend.Listen()

	go home.OpenMainWindow()
	app.Main()
	//UIInit()
	//backend.Listen()
}
