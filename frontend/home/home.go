package home

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/JesperY/GoCopyUSB/config"
	page "github.com/JesperY/GoCopyUSB/frontend/pages"
	"github.com/JesperY/GoCopyUSB/frontend/pages/settings"
	"log"
)

func OpenMainWindow() {
	window := new(app.Window)
	window.Option(app.Title(config.ConfigPtr.Title), app.Size(unit.Dp(config.ConfigPtr.Width), unit.Dp(config.ConfigPtr.Height)))
	config.ConfigPtr.Win = window
	fmt.Println("当前win的值是", config.ConfigPtr.Win == nil)
	if err := loop(window); err != nil {
		log.Fatal(err)
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
			//window.Perform(system.ActionClose)
			fmt.Println("按下X按钮")
			config.ConfigPtr.Win = nil
			fmt.Println("当前win的值是", config.ConfigPtr.Win == nil)
			return nil
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			router.Layout(gtx, th)
			e.Frame(gtx.Ops)
		}
	}
}
