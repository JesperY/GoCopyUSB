package settings

import (
	"gioui.org/app"
	"gioui.org/example/component/icon"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/JesperY/GoCopyUSB/config"
	alo "github.com/JesperY/GoCopyUSB/frontend/applayout"
	page "github.com/JesperY/GoCopyUSB/frontend/pages"
	"github.com/gen2brain/dlgs"
	"log"
)

type Page struct {
	chooseDirBtn widget.Clickable
	app.Window
	*page.Router
	widget.List
}

func (p Page) Actions() []component.AppBarAction {
	return nil
}

func (p Page) Overflow() []component.OverflowAction {
	return nil
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
		return layout.Flex{
			Alignment: layout.Middle,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DefaultInset.Layout(gtx, material.Body1(th, `感谢使用USBCopier！`).Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DetailRow{}.Layout(gtx, material.Body1(th, "当前备份位置").Layout, func(gtx layout.Context) layout.Dimensions {
					if p.chooseDirBtn.Clicked(gtx) {
						path, _, err := dlgs.File("Select Folder", "", true)
						if err != nil {
							log.Println("Error selecting folder:", err)
						} else {
							config.ConfigPtr.TargetDir = path
							log.Println(path)
							p.Window.Invalidate()
						}
					}
					return material.Button(th, &p.chooseDirBtn, config.ConfigPtr.TargetDir).Layout(gtx)
				})
			}),
		)
	})
}

func (p *Page) NavItem() component.NavItem {
	return component.NavItem{
		Name: "参数配置",
		Icon: icon.HomeIcon,
	}
}

func New(router *page.Router) *Page {
	return &Page{
		Router: router,
	}
}
