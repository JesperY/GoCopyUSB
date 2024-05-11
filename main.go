package main

import (
	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/JesperY/GoCopyUSB/backEnd"
	"image/color"
	"os"
)

func main() {
	//UIinit()
	backEnd.Listen()
}

func UIinit() {
	var title string = `USBCopier`
	var width, height int = 600, 400
	go func() {
		window := new(app.Window)
		SetWindowOptions(window, title, width, height)
		ListenEvent(window)
	}()
	app.Main()

}

// SetWindowOptions
//
//	@Description: 设置窗口的各项属性
//	@param title 标题文字
//	@param width 窗口长度 单位DP 下同
//	@param height 窗口宽度
func SetWindowOptions(window *app.Window, title string, width, height int) {
	window.Option(app.Title(title), app.Size(unit.Dp(width), unit.Dp(height)))
}

// ListenEvent
//
//	@Description: 监听页面内容
//	@param window 窗口对象指针
func ListenEvent(window *app.Window) {
	theme := material.NewTheme()
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			// Define an large label with an appropriate text:
			title := material.H1(theme, "Init success")

			// Change the color of the label.
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			title.Color = maroon

			// Change the position of the label.
			title.Alignment = text.Middle

			// Draw the label to the graphics context.
			title.Layout(gtx)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
