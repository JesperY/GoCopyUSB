package main

import (
	"gioui.org/app"
	"github.com/JesperY/GoCopyUSB/backend"
	"github.com/JesperY/GoCopyUSB/copylogger"
	"github.com/JesperY/GoCopyUSB/frontend/home"
	"github.com/JesperY/GoCopyUSB/frontend/tray"
)

func main() {
	defer copylogger.SugarLogger.Sync()
	go tray.SysTrayRun()
	go backend.Listen()
	go home.OpenMainWindow()
	app.Main()
	//UIInit()
	//backend.Listen()
}
