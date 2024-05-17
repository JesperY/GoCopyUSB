package tray

import (
	"fmt"
	"github.com/JesperY/GoCopyUSB/config"
	"github.com/JesperY/GoCopyUSB/frontend/home"
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"os"
	"time"
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
		for {
			select {
			case <-time.Tick(time.Second):
				fmt.Println("还在监听...")
			case <-mReOpen.ClickedCh:
				fmt.Println("按下打开按钮")
				if config.ConfigPtr.Win == nil {
					home.OpenMainWindow()
				}
			}
		}
	}()
	mQuit := systray.AddMenuItem("退出", "退出")
	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}
func OnExit() {
	systray.Quit()
	os.Exit(0)
}
