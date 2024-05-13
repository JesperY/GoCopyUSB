package pages

import (
	"gioui.org/example/component/icon"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"log"
	"time"
)

type Page interface {
	Actions() []component.AppBarAction
	Overflow() []component.OverflowAction
	Layout(gtx layout.Context, th *material.Theme) layout.Dimensions
	NavItem() component.NavItem
}

type Router struct {
	pages   map[interface{}]Page
	current interface{}
	*component.ModalNavDrawer
	NavAnim component.VisibilityAnimation
	*component.AppBar
	*component.ModalLayer
	NonModalDrawer, BottomBar bool
}

func NewRouter() Router {
	modal := component.NewModal()

	nav := component.NewNav("导航", "USB导航")
	modalNav := component.ModalNavFrom(&nav, modal)

	bar := component.NewAppBar(modal)
	bar.NavigationIcon = icon.MenuIcon

	na := component.VisibilityAnimation{
		State:    component.Invisible,
		Duration: time.Millisecond * 200,
	}
	return Router{
		pages:          make(map[interface{}]Page),
		ModalLayer:     modal,
		ModalNavDrawer: modalNav,
		AppBar:         bar,
		NavAnim:        na,
	}
}

// Register  注册路由
func (r *Router) Register(tag interface{}, p Page) {
	r.pages[tag] = p
	navItem := p.NavItem()
	navItem.Tag = tag
	if r.current == interface{}(nil) {
		r.current = tag
		r.AppBar.Title = navItem.Name
		r.AppBar.SetActions(p.Actions(), p.Overflow())
	}
	r.ModalNavDrawer.AddNavItem(navItem)
}

// SwitchTo 切换标签页
func (r *Router) SwitchTo(tag interface{}) {
	p, ok := r.pages[tag]
	if !ok {
		return
	}
	navItem := p.NavItem()
	r.current = tag
	r.AppBar.Title = navItem.Name
	r.AppBar.SetActions(p.Actions(), p.Overflow())
}

// 布局
func (r *Router) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for _, event := range r.AppBar.Events(gtx) {
		switch e := event.(type) {
		case component.AppBarNavigationClicked:
			if r.NonModalDrawer {
				r.NavAnim.ToggleVisibility(gtx.Now)
			} else {
				r.ModalNavDrawer.Appear(gtx.Now)
				r.NavAnim.Disappear(gtx.Now)
			}
		case component.AppBarContextMenuDismissed:
			log.Printf("Context menu dismissed: %v", e)
		case component.AppBarOverflowActionClicked:
			log.Printf("Overflow action selected: %v", e)
		}
	}
	if r.ModalNavDrawer.NavDestinationChanged() {
		r.SwitchTo(r.ModalNavDrawer.CurrentNavDestination())
	}
	paint.Fill(gtx.Ops, th.Palette.Bg)

	content := layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X /= 3
				return r.NavDrawer.Layout(gtx, th, &r.NavAnim)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return r.pages[r.current].Layout(gtx, th)
			}),
		)
	})
	bar := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return r.AppBar.Layout(gtx, th, "Menu", "Actions")
	})
	flex := layout.Flex{Axis: layout.Vertical}
	if r.BottomBar {
		flex.Layout(gtx, content, bar)
	} else {
		flex.Layout(gtx, bar, content)
	}
	r.ModalLayer.Layout(gtx, th)
	return layout.Dimensions{Size: gtx.Constraints.Max}
}