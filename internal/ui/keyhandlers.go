package ui

import (
	"container/ring"
	"fmt"
	"github.com/CaninoDev/go-hackernews/internal/api"
	"github.com/gdamore/tcell/v2"
	"time"
)


func delaySeconds(n time.Duration, d *Display) {
	endpoints := d.DB.EndPoints()
	ring := ring.New(len(endpoints))

	for _, endpoint := range endpoints {
		ring.Value = endpoint
		ring.Next()
	}

	for _ = range time.Tick(n * time.Second) {
		time.Sleep(2 * time.Second)
		d.DebugBar.Clear()
		d.DebugBar.SetText(fmt.Sprintf("%v", ring.Value))
	}

}
func (d *Display) AppKeyHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEscape:
		d.App.Stop()
	case tcell.KeyCtrlN:
		d.Posts.SetCurrentTab(api.NewS.String())
	case tcell.KeyCtrlJ:
		d.Posts.SetCurrentTab(api.Jobs.String())
	case tcell.KeyCtrlT:
		d.Posts.SetCurrentTab(api.Top.String())
	case tcell.KeyCtrlB:
		d.Posts.SetCurrentTab(api.Best.String())
	case tcell.KeyCtrlS:
		d.Posts.SetCurrentTab(api.Show.String())
	case tcell.KeyCtrlA:
		d.Posts.SetCurrentTab(api.Ask.String())
	case tcell.KeyCtrlG:
		d.generatePostList()
	case tcell.KeyCtrlC:
		d.App.SetFocus(d.Comments.Tree)
	}
	return event
}
