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
//			//a.listView.debugBar.SetText(endpoint.String())
//			a.listView.tabbedLists.SetCurrentTab(endpoint.String())
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

	cycleFocus := func() {
		if panel == POSTPANEL {
			if a.postView.content.HasFocus() {
				a.ui.SetFocus(a.postView.commentsTree)
			} else {
				a.panels.SetCurrentPanel(LISTPANEL)
			}
		} else {
			a.panels.SetCurrentPanel(POSTPANEL)
			a.ui.SetFocus(a.postView.content)
		}
	}

	a.ui.QueueUpdateDraw(cycleFocus)

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
		a.quit(event)
	case tcell.KeyTab:
		a.handleToggle()
	case tcell.KeyCtrlN:
		a.listView.tabbedLists.SetCurrentTab(api.New.String())
	case tcell.KeyCtrlJ:
		a.listView.tabbedLists.SetCurrentTab(api.Jobs.String())
	case tcell.KeyCtrlT:
		a.listView.tabbedLists.SetCurrentTab(api.Top.String())
	case tcell.KeyCtrlB:
		a.listView.tabbedLists.SetCurrentTab(api.Best.String())
	case tcell.KeyCtrlS:
		a.listView.tabbedLists.SetCurrentTab(api.Show.String())
	case tcell.KeyCtrlA:
		a.listView.tabbedLists.SetCurrentTab(api.Ask.String())
	case tcell.KeyCtrlP:
		a.listView.pageNav(prev)
	case tcell.KeyCtrlL:
		a.listView.pageNav(next)
	default:
		return event
	}
	go a.listView.populateList()
	return event
}
