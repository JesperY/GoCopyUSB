package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/JesperY/GoCopyUSB/backend"
	"github.com/JesperY/GoCopyUSB/copylogger"
	page "github.com/JesperY/GoCopyUSB/frontend/pages"
	"github.com/JesperY/GoCopyUSB/frontend/pages/settings"
	"os"
)

func main() {
	defer copylogger.SugarLogger.Sync()
	go backend.Listen()
	go func() {
		window := new(app.Window)
		if err := loop(window); err != nil {
		}
		os.Exit(0)
	}()
	app.Main()
	//UIInit()
	//backend.Listen()
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
