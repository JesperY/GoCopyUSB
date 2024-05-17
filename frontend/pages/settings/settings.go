package settings

import (
	"context"
	"fmt"
	"gioui.org/app"
	"gioui.org/example/component/icon"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
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
	"strconv"
	"strings"
	"sync"
)

type Page struct {
	chooseDirBtn   widget.Clickable
	chooseBlkBtn   widget.Clickable
	clearBlkBtn    widget.Clickable
	submitSufBtn   widget.Clickable
	submitDelayBtn widget.Clickable
	app.Window
	*page.Router
	widget.List
	textList          widget.List
	text              []string
	blkSuffix         component.TextField
	delayTime         component.TextField
	inputAlignment    text.Alignment
	inputAlignment2   text.Alignment
	dialogOpen        bool
	closeDialogBtn    widget.Clickable
	App               *Application
	autoLaunchTrigger widget.Bool
}

func (p Page) Actions() []component.AppBarAction {
	return nil
}

func (p Page) Overflow() []component.OverflowAction {
	return nil
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	p.List.Axis = layout.Vertical
	ctx := context.Background()
	a := NewApplication(ctx)
	p.App = a
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
											config.ConfigPtr.WriteConfig()
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
											config.ConfigPtr.WhiteListDir = nil
											config.ConfigPtr.WriteConfig()
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

						border := widget.Border{
							Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
							CornerRadius: unit.Dp(4),
							Width:        unit.Dp(2),
						}
						return layout.Inset{Left: unit.Dp(75), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								// 使用 layout.Stack 来确保内容和滚动条正确布局
								return layout.Stack{}.Layout(gtx, layout.Expanded(func(gtx layout.Context) layout.Dimensions {
									p.textList.Axis = layout.Vertical
									// 设置文本框的大小
									textboxWidth := gtx.Constraints.Max.X
									textboxHeight := gtx.Dp(unit.Dp(150))
									gtx.Constraints.Min.X = textboxWidth
									gtx.Constraints.Min.Y = textboxHeight
									gtx.Constraints.Max.Y = textboxHeight
									return material.List(th, &p.textList).Layout(gtx, len(config.ConfigPtr.WhiteListDir), func(gtx layout.Context, index int) layout.Dimensions {
										return layout.Stack{}.Layout(gtx,
											layout.Stacked(func(gtx layout.Context) layout.Dimensions {
												// 设置行高
												gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(30))
												gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(30))
												gtx.Constraints.Min.X = textboxWidth - gtx.Dp(unit.Dp(2))*2 // 减去边框的宽度
												gtx.Constraints.Max.X = textboxWidth - gtx.Dp(unit.Dp(2))*2
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
						})
					}),
				)
			}),
			// 添加一些空格
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				spacerWidth := unit.Dp(20)
				return layout.Spacer{Width: spacerWidth, Height: unit.Dp(10)}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{Size: gtx.Constraints.Min}
						})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						p.blkSuffix.Alignment = p.inputAlignment
						return p.blkSuffix.Layout(gtx, th, `请在此输入您不想备份的文件后缀名，以空格隔开`)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if p.submitSufBtn.Clicked(gtx) {
							str := p.blkSuffix.Text()
							split := strings.Split(str, " ")
							for _, s := range split {
								for _, suffix := range config.ConfigPtr.WhiteListSuffix {
									if s != suffix {
										config.ConfigPtr.WhiteListSuffix = append(config.ConfigPtr.WhiteListSuffix, s)
									}
								}
							}
							config.ConfigPtr.WriteConfig()
						}
						// 设置按钮的最小和最大宽度
						btn := material.Button(th, &p.submitSufBtn, "确认")
						//btnSize := layout.Dimensions{Size: gtx.Constraints.Min}
						return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(8), Top: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Dp(80)
							gtx.Constraints.Max.X = gtx.Dp(80)
							return btn.Layout(gtx)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				var suffixStr string
				for _, suffix := range config.ConfigPtr.WhiteListSuffix {
					suffixStr += suffix + ","
				}
				return alo.DefaultInset.Layout(gtx, material.Body2(th, fmt.Sprintf(`当前黑名单后缀：%s`, suffixStr)).Layout)
			}),
			// 添加一些空格
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				spacerWidth := unit.Dp(20)
				return layout.Spacer{Width: spacerWidth, Height: unit.Dp(10)}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Horizontal,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{Size: gtx.Constraints.Min}
						})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						p.delayTime.Alignment = p.inputAlignment2
						return p.delayTime.Layout(gtx, th, `请在此输入要延迟备份的分钟数(请输入整数)`)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if p.submitDelayBtn.Clicked(gtx) {
							str := p.delayTime.Text()
							minute, err := stringToInt(str)
							if err != nil {
								p.delayTime.SetText("")
								p.showErrDialog("err", "请输入一个有效的整数")
							} else {
								config.ConfigPtr.DelayMinutes = minute
								config.ConfigPtr.WriteConfig()
								fmt.Println(minute)
							}
						}
						// 设置按钮的最小和最大宽度
						btn := material.Button(th, &p.submitDelayBtn, "确认")
						//btnSize := layout.Dimensions{Size: gtx.Constraints.Min}
						return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(8), Top: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Dp(80)
							gtx.Constraints.Max.X = gtx.Dp(80)
							return btn.Layout(gtx)
						})
					}),
				)
			}),
			// 添加一些空格
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				spacerWidth := unit.Dp(20)
				return layout.Spacer{Width: spacerWidth, Height: unit.Dp(10)}.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DetailRow{}.Layout(gtx,
					material.Body1(th, "是否开机自启动").Layout,
					func(gtx layout.Context) layout.Dimensions {
						if p.autoLaunchTrigger.Update(gtx) {
							if p.autoLaunchTrigger.Value {
								config.ConfigPtr.AutoStartUp = true
							} else {
								config.ConfigPtr.AutoStartUp = true
							}
							config.ConfigPtr.WriteConfig()
						}
						return material.Switch(th, &p.autoLaunchTrigger, "是否开机自启动").Layout(gtx)
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
func stringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// Application keeps track of all the windows and global state.
type Application struct {
	// Context is used to broadcast application shutdown.
	Context context.Context
	// Shutdown shuts down all windows.
	Shutdown func()
	// active keeps track the open windows, such that application
	// can shut down, when all of them are closed.
	active sync.WaitGroup
}

// NewWindow creates a new tracked window.
func (a *Application) NewWindow(title string, view View, opts ...app.Option) {
	opts = append(opts, app.Title(title), app.Size(unit.Dp(200), unit.Dp(100)))

	a.active.Add(1)
	go func() {
		defer a.active.Done()

		w := &Window{
			App:    a,
			Window: new(app.Window),
		}
		w.Window.Option(opts...)
		view.Run(w)
	}()
}

type View interface {
	// Run handles the window event loop.
	Run(w *Window) error
}
type Window struct {
	App *Application
	*app.Window
}

// WidgetView allows to use layout.Widget as a view.
type WidgetView func(gtx layout.Context, th *material.Theme) layout.Dimensions

// Run displays the widget with default handling.
func (view WidgetView) Run(w *Window) error {
	var ops op.Ops

	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	go func() {
		<-w.App.Context.Done()
		w.Perform(system.ActionClose)
	}()
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			view(gtx, th)
			e.Frame(gtx.Ops)
		}
	}
}
func NewApplication(ctx context.Context) *Application {
	ctx, cancel := context.WithCancel(ctx)
	return &Application{
		Context:  ctx,
		Shutdown: cancel,
	}
}

// Wait waits for all windows to close.
func (a *Application) Wait() {
	a.active.Wait()
}

// showErrDialog
//
//	@Description: 显示错误弹窗
//	@receiver p   页面指针
//	@param title  窗口标题
//	@param info   提示内容
func (p *Page) showErrDialog(title, info string) {
	p.App.NewWindow(title,
		WidgetView(func(gtx layout.Context, th *material.Theme) layout.Dimensions {

			return layout.Center.Layout(gtx, material.Body2(th, info).Layout)
		}),
		//app.Size(size, size),
	)
}
