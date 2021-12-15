package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/CaninoDev/go-hackernews/internal/api"
	"github.com/gdamore/tcell/v2"
)

//func (a *App) initializeBindings() {
//	bindingCfg := cbind.NewConfiguration()
//
//	bindingCfg.SetKey(tcell.ModNone, tcell.KeyTab, a.handleToggle)
//
//	bindingCfg.SetKey(tcell.ModNone, tcell.KeyEscape, a.quit)
//
//	bindingCfg.SetKey(tcell.ModCtrl, 'n', a.handleTabSwitch)
//	bindingCfg.Set("Ctrl+n", )
//	for _, endpoint := range api.AllEndPoints() {
//		endpointRune := []rune(strings.ToLower(endpoint.String()[0:1]))[0]
//		bindingCfg.SetRune(tcell.ModNone, endpointRune, func(event *tcell.EventKey) *tcell.EventKey {
//			//a.browser.debugBar.SetText(endpoint.String())
//			a.browser.lists.SetCurrentTab(endpoint.String())
//			return nil
//		})
//	}
//	a.ui.SetInputCapture(bindingCfg.Capture)
//}

func (a *App) handleTabSwitch(ev *tcell.EventKey) *tcell.EventKey {
	return nil
}

func (a *App) quit(_ *tcell.EventKey) *tcell.EventKey {
	a.ui.Stop()
	return nil
}

func (a *App) handleToggle() *tcell.EventKey {
	panel, _ := a.panels.GetFrontPanel()

	a.browser.debugBar.SetText(panel)
	if panel == "post" {
		a.panels.SetCurrentPanel("browser")
	} else {
		a.panels.SetCurrentPanel("post")
	}
	return nil
}

func (a *App) initializeInputHandler(panes ...cview.Primitive) {
	panes = append(panes)
	a.focusManager = cview.NewFocusManager(a.ui.SetFocus)
	a.focusManager.SetWrapAround(true)
	a.focusManager.Add(panes...)
}

func (a *App) inputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		a.ui.Stop()
	case tcell.KeyTab:
		a.handleToggle()
	case tcell.KeyCtrlN:
		a.browser.currentTab = api.New.String()
		a.browser.lists.SetCurrentTab(api.New.String())
	case tcell.KeyCtrlJ:
		a.browser.currentTab = api.Jobs.String()
		a.browser.lists.SetCurrentTab(api.Jobs.String())
	case tcell.KeyCtrlT:
		a.browser.currentTab = api.Top.String()
		a.browser.lists.SetCurrentTab(api.Top.String())
	case tcell.KeyCtrlB:
		a.browser.currentTab = api.Best.String()
		a.browser.lists.SetCurrentTab(api.Best.String())
	case tcell.KeyCtrlS:
		a.browser.currentTab = api.Show.String()
		a.browser.lists.SetCurrentTab(api.Show.String())
	case tcell.KeyCtrlA:
		a.browser.currentTab = api.Ask.String()
		a.browser.lists.SetCurrentTab(api.Ask.String())
	case tcell.KeyCtrlP:
		a.browser.pageNav(false)
	case tcell.KeyCtrlL:
		a.browser.pageNav(true)
	default:
		return event
	}
	go a.browser.populateList()
	return event
}
