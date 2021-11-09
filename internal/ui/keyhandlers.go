package ui

import (
	"github.com/CaninoDev/go-hackernews/internal/api"
	"github.com/gdamore/tcell/v2"
)

func (d *Display) AppKeyHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		d.App.Stop()
	case tcell.KeyCtrlN:
		d.List(api.NewS)
		d.App.SetFocus(d.Posts)
	case tcell.KeyCtrlJ:
		d.List(api.Jobs)
		d.App.SetFocus(d.Posts)
	case tcell.KeyCtrlT:
		d.List(api.Top)
		d.App.SetFocus(d.Posts)
	case tcell.KeyCtrlB:
		d.List(api.Best)
		d.App.SetFocus(d.Posts)
	case tcell.KeyCtrlS:
		d.List(api.Show)
		d.App.SetFocus(d.Posts)
	case tcell.KeyCtrlA:
		d.List(api.Ask)
		d.App.SetFocus(d.Posts)
	case tcell.KeyCtrlG:
		d.App.SetFocus(d.Comments.Tree)
	}
	return event
}

