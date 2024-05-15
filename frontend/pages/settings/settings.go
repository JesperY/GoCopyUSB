package settings

import (
	"gioui.org/app"
	"gioui.org/example/component/icon"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/JesperY/GoCopyUSB/config"
	alo "github.com/JesperY/GoCopyUSB/frontend/applayout"
	page "github.com/JesperY/GoCopyUSB/frontend/pages"
	"github.com/gen2brain/dlgs"
	"github.com/ncruces/zenity"
	"image/color"
	"log"
)

type Page struct {
	chooseDirBtn widget.Clickable
	chooseBlkBtn widget.Clickable
	clearBlkBtn  widget.Clickable
	setDefault   widget.Bool
	app.Window
	*page.Router
	widget.List
	textList widget.List
	text     []string
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
			// 介绍文本
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DefaultInset.Layout(gtx, material.Body1(th, `感谢使用USBCopier！`).Layout)
			}),
			// 当前备份位置和选择按钮
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DetailRow{}.Layout(gtx, material.Body1(th, "当前备份位置").Layout, func(gtx layout.Context) layout.Dimensions {
					if p.chooseDirBtn.Clicked(gtx) {
						path, _, err := dlgs.File("Select Folder", "", true)
						if err != nil {
							log.Println("Error selecting folder:", err)
						} else {
							if path == "" {
								_, err := dlgs.Info("路径不可用", "请选择一个有效的文件夹路径。")
								if err != nil {
									log.Println("Error displaying dialog:", err)
								}
							}
							config.ConfigPtr.TargetDir = path
							config.ConfigPtr.WriteConfig()
							p.Window.Invalidate()
						}
					}
					return material.Button(th, &p.chooseDirBtn, config.ConfigPtr.TargetDir).Layout(gtx)
				})
			}),
			// 水平排列黑名单按钮和文本框
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					// 垂直排列两个按钮
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(gtx,
							// 第一个黑名单按钮
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								// 固定按钮宽度和高度
								buttonWidth := unit.Dp(150)
								buttonHeight := unit.Dp(50)
								// 外边距
								inset := layout.Inset{
									Top:    unit.Dp(10),
									Bottom: unit.Dp(10),
									Left:   unit.Dp(10),
									Right:  unit.Dp(10),
								}
								return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return layout.Stack{}.Layout(gtx, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
										// 设置按钮的最小宽度和高度
										gtx.Constraints.Min.X = gtx.Dp(buttonWidth)
										gtx.Constraints.Min.Y = gtx.Dp(buttonHeight)
										// 设置按钮的最大宽度和高度（可选）
										gtx.Constraints.Max.X = gtx.Dp(buttonWidth)
										gtx.Constraints.Max.Y = gtx.Dp(buttonHeight)
										if p.chooseBlkBtn.Clicked(gtx) {
											folders, err := zenity.SelectFileMultiple(
												zenity.Title("选择不需要备份的文件夹"),
												zenity.Directory(),
											)
											if err != nil {
												log.Println("Error selecting folder:", err)
											}
											for _, folder := range folders {
												config.ConfigPtr.WhiteListDir = append(config.ConfigPtr.WhiteListDir, folder)
											}
											p.Window.Invalidate()
										}
										return material.Button(th, &p.chooseBlkBtn, "点击选择黑名单").Layout(gtx)
									}))
								})
							}),
							// 在两个按钮之间添加空格
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								spacerWidth := unit.Dp(20)
								return layout.Spacer{Width: spacerWidth, Height: unit.Dp(10)}.Layout(gtx)
							}),
							// 第二个黑名单按钮
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								// 固定按钮宽度和高度
								buttonWidth := unit.Dp(150)
								buttonHeight := unit.Dp(50)
								// 外边距
								inset := layout.Inset{
									Top:    unit.Dp(10),
									Bottom: unit.Dp(10),
									Left:   unit.Dp(10),
									Right:  unit.Dp(10),
								}
								return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return layout.Stack{}.Layout(gtx, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
										// 设置按钮的最小宽度和高度
										gtx.Constraints.Min.X = gtx.Dp(buttonWidth)
										gtx.Constraints.Min.Y = gtx.Dp(buttonHeight)
										// 设置按钮的最大宽度和高度（可选）
										gtx.Constraints.Max.X = gtx.Dp(buttonWidth)
										gtx.Constraints.Max.Y = gtx.Dp(buttonHeight)
										if p.clearBlkBtn.Clicked(gtx) {
											p.Window.Invalidate()
										}
										return material.Button(th, &p.clearBlkBtn, "清空黑名单").Layout(gtx)
									}))
								})
							}),
						)
					}),
					// 在右侧添加一个带边框且固定大小的可滚动文本框
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						// 设置文本框的大小
						textboxWidth := gtx.Dp(unit.Dp(400))
						textboxHeight := gtx.Dp(unit.Dp(150))
						gtx.Constraints.Min.X = textboxWidth
						gtx.Constraints.Min.Y = textboxHeight
						gtx.Constraints.Max.X = textboxWidth
						gtx.Constraints.Max.Y = textboxHeight

						border := widget.Border{
							Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
							CornerRadius: unit.Dp(4),
							Width:        unit.Dp(2),
						}
						return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							// 使用 layout.Stack 来确保内容和滚动条正确布局
							return layout.Stack{}.Layout(gtx, layout.Expanded(func(gtx layout.Context) layout.Dimensions {
								p.textList.Axis = layout.Vertical
								return material.List(th, &p.textList).Layout(gtx, len(config.ConfigPtr.WhiteListDir), func(gtx layout.Context, index int) layout.Dimensions {
									// 设置背景颜色
									bgColor := color.NRGBA{R: 240, G: 240, B: 240, A: 255} // 默认颜色
									if index%2 == 0 {
										bgColor = color.NRGBA{R: 220, G: 220, B: 220, A: 255} // 另一种颜色
									}

									return layout.Stack{}.Layout(gtx,
										layout.Expanded(func(gtx layout.Context) layout.Dimensions {
											paint.Fill(gtx.Ops, bgColor)
											gtx.Constraints.Min.Y = 30
											gtx.Constraints.Max.X = 400
											return layout.Dimensions{Size: gtx.Constraints.Min}
										}),
										layout.Stacked(func(gtx layout.Context) layout.Dimensions {
											// 设置行高
											gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(30))
											gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(30))
											gtx.Constraints.Min.X = gtx.Dp(unit.Dp(400))
											gtx.Constraints.Max.X = gtx.Dp(unit.Dp(400))
											insets := layout.UniformInset(unit.Dp(0))
											return insets.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													return material.Body1(th, config.ConfigPtr.WhiteListDir[index]).Layout(gtx)
												})
											})
										}),
									)
								})
							}))
						})
					}),
				)
			}),
			//// 黑名单
			//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			//	// 固定按钮宽度和高度
			//	buttonWidth := unit.Dp(150)
			//	buttonHeight := unit.Dp(50)
			//	// 外边距
			//	inset := layout.Inset{
			//		Left: unit.Dp(10),
			//	}
			//	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			//		return layout.Stack{}.Layout(gtx, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			//			// 设置按钮的最小宽度和高度
			//			gtx.Constraints.Min.X = gtx.Dp(buttonWidth)
			//			gtx.Constraints.Min.Y = gtx.Dp(buttonHeight)
			//			// 设置按钮的最大宽度和高度（可选）
			//			gtx.Constraints.Max.X = gtx.Dp(buttonWidth)
			//			gtx.Constraints.Max.Y = gtx.Dp(buttonHeight)
			//			if p.chooseDirBtn.Clicked(gtx) {
			//				path, _, err := dlgs.File("Select Folder", "", true)
			//				if err != nil {
			//					log.Println("Error selecting folder:", err)
			//				} else {
			//					if path == "" {
			//						_, err := dlgs.Info("路径不可用", "请选择一个有效的文件夹路径。")
			//						if err != nil {
			//							log.Println("Error displaying dialog:", err)
			//						}
			//					}
			//					config.ConfigPtr.TargetDir = path
			//					config.ConfigPtr.WriteConfig()
			//					p.Window.Invalidate()
			//				}
			//			}
			//			return material.Button(th, &p.chooseBlkBtn, "点击选择黑名单").Layout(gtx)
			//		}))
			//
			//	})
			//}),
			//// 在两个按钮之间添加空格
			//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			//	spacerWidth := unit.Dp(20)
			//	return layout.Spacer{Width: spacerWidth, Height: unit.Dp(10)}.Layout(gtx)
			//}),
			//layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			//	// 固定按钮宽度和高度
			//	buttonWidth := unit.Dp(150)
			//	buttonHeight := unit.Dp(50)
			//	// 外边距
			//	inset := layout.Inset{
			//		Left: unit.Dp(10),
			//	}
			//	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			//		return layout.Stack{}.Layout(gtx, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			//			// 设置按钮的最小宽度和高度
			//			gtx.Constraints.Min.X = gtx.Dp(buttonWidth)
			//			gtx.Constraints.Min.Y = gtx.Dp(buttonHeight)
			//			// 设置按钮的最大宽度和高度（可选）
			//			gtx.Constraints.Max.X = gtx.Dp(buttonWidth)
			//			gtx.Constraints.Max.Y = gtx.Dp(buttonHeight)
			//			if p.chooseDirBtn.Clicked(gtx) {
			//				path, _, err := dlgs.File("Select Folder", "", true)
			//				if err != nil {
			//					log.Println("Error selecting folder:", err)
			//				} else {
			//					if path == "" {
			//						_, err := dlgs.Info("路径不可用", "请选择一个有效的文件夹路径。")
			//						if err != nil {
			//							log.Println("Error displaying dialog:", err)
			//						}
			//					}
			//					config.ConfigPtr.TargetDir = path
			//					config.ConfigPtr.WriteConfig()
			//					p.Window.Invalidate()
			//				}
			//			}
			//			return material.Button(th, &p.chooseBlkBtn, "点击选择黑名单").Layout(gtx)
			//		}))
			//	})
			//}),
			////TODO  在下方添加代码，绘制一个右侧的文本框
			//// 在右侧添加一个可滚动的文本框
			//layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			//	return layout.Stack{}.Layout(gtx, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			//		// 定义文本框内容
			//		text := []string{"Line 1", "Line 2", "Line 3", "Line 4", "Line 5", "Line 6", "Line 7", "Line 8"}
			//		list := layout.List{Axis: layout.Vertical}
			//		return list.Layout(gtx, len(text), func(gtx layout.Context, index int) layout.Dimensions {
			//			return material.Body1(th, text[index]).Layout(gtx)
			//		})
			//	}))
			//}),
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
