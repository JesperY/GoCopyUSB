package tray

import (
	"github.com/JesperY/GoCopyUSB/config"
	"github.com/JesperY/GoCopyUSB/frontend/home"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"os"
)

func SysTrayRun() {
	systray.Run(OnReady, OnExit)
}

func OnReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("USBCopier")
	systray.SetTooltip("A GODDAMN SOFTWARE!")
	mReOpen := systray.AddMenuItem("打开主面板", "打开主面板")
	go func() {
		<-mReOpen.ClickedCh
		if config.ConfigPtr.Win == nil {
			home.OpenMainWindow()
		}
	}()
	mQuit := systray.AddMenuItem("退出", "退出")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
		os.Exit(0)
	}()
}
func OnExit() {
	systray.Quit()
	os.Exit(0)
}
