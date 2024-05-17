package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/JesperY/GoCopyUSB/backend"
	"github.com/JesperY/GoCopyUSB/config"
	"github.com/JesperY/GoCopyUSB/copylogger"
	page "github.com/JesperY/GoCopyUSB/frontend/pages"
	"github.com/JesperY/GoCopyUSB/frontend/pages/settings"
	"os"
)

func main() {
	initCheck()
	defer copylogger.SugarLogger.Sync()
	go backend.Listen()
	go func() {
		window := new(app.Window)
		if err := loop(window); err != nil {
		}
		//os.Exit(0)
	}()
	app.Main()
	//UIInit()
	//backend.Listen()
}

func initCheck() {
	// todo 目标目录检查
	targetPath := config.ConfigPtr.TargetDir
	fileInfo, err := os.Stat(targetPath)
	// 目标文件夹不存在，创建
	if err != nil {
		err := os.Mkdir(targetPath, 0777)
		if err != nil {
			// todo 没有创建权限，无法继续，弹窗
			copylogger.SugarLogger.Errorf("Permission deny, can not create target dir, %v", err)
			return
		}
	}
	if fileInfo.Mode().Perm()&os.ModePerm == 0 {
		// todo 没有修改权限，无法继续，弹窗
		copylogger.SugarLogger.Errorf("Permission deny, can not modify target dir")
	}

	// todo 开机自启动项检查？
	if config.ConfigPtr.AutoStartUp {
		// 配置开机自启动，检查注册表项
		// 如果未配置自启动则配置
		if !backend.IsAutoStartUp() {
			// 已配置自启动
			//backend.DisableAutoStartUp()
			backend.EnableAutoStartUp()
		}

	} else {
		// 配置开机不启动，检查注册表项
		// 如果已配置则删除
		if backend.IsAutoStartUp() {
			backend.DisableAutoStartUp()
		}

	}

}

func loop(window *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	// 注册路由
	router := page.NewRouter()
	router.Register(0, settings.New(&router)) // 设置属性页面
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			router.Layout(gtx, th)
			e.Frame(gtx.Ops)
		}
	}

}
