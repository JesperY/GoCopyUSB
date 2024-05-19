package main

import (
	"gioui.org/app"
	"github.com/JesperY/GoCopyUSB/backend"
	"github.com/JesperY/GoCopyUSB/frontend/home"
	"github.com/JesperY/GoCopyUSB/frontend/tray"
	"github.com/JesperY/GoCopyUSB/logger"
)

func main() {
	// todo 重复运行检查
	home.ErrorDialog("test", "errorTest")
	//page := &settings.Page{}
	//page.ShowErrDialog("test", "error test")
	fileLock := backend.SingleCheck()
	if fileLock != nil {
		defer fileLock.Unlock()
	}
	backend.InitCheck()
	defer logger.SugarLogger.Sync()
	go tray.SysTrayRun()
	go backend.Listen()

	go home.OpenMainWindow()
	app.Main()
	//UIInit()
	//backend.Listen()
}
