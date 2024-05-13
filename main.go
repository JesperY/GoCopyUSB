package test

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/JesperY/GoCopyUSB/backEnd"
	"github.com/JesperY/GoCopyUSB/config"
	"image/color"
)

func main() {
	UIInit()
}

func UIInit() {
	go func() {
		window := new(app.Window)
		SetWindowOptions(window, config.Title, config.Width, config.Height)
		ListenEvent(window)
	}()
	go backEnd.Listen()
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
			//os.Exit(0)
			//todo 点击关闭时程序缩小至托盘
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			// Define an large label with an appropriate text:
			title := material.H6(theme, "USB备份监听已启动")
			dirText := material.Body1(theme, fmt.Sprintf(`当前备份目标目录为%s`, config.TargetDir))

			// Change the color of the label.
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			title.Color = maroon
			dirColor := color.NRGBA{R: 127, G: 55, B: 123, A: 255}
			dirText.Color = dirColor
			// Change the position of the label.
			title.Alignment = text.Middle
			dirText.Alignment = text.Middle
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(title.Layout),
				layout.Rigid(dirText.Layout),
			)
			e.Frame(gtx.Ops)
		}
	}
}
