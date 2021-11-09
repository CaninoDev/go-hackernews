package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"fmt"
	"github.com/CaninoDev/go-hackernews/internal/api"
	"log"
	"net/url"
)

type Posts struct {
	*cview.TabbedPanels
}

func (d *Display) NewPosts() *Posts {
	endPoints := []api.EndPoint{
		api.NewS, api.Top, api.Best, api.Ask, api.Show,
	}
	d.Listings = make(map[string]*cview.List)
	tabbedPanels := cview.NewTabbedPanels()
	for _, endpoint := range endPoints {
		d.Listings[endpoint.String()] = d.generatePostsList(endpoint)
		tabbedPanels.AddTab(endpoint.String(), endpoint.String(), d.Listings[endpoint.String()])
	}
	return &Posts{
		tabbedPanels,
	}
}

func (d *Display) List(endPoint api.EndPoint) {
	d.Listings[endPoint.String()] = d.generatePostsList(endPoint)
	d.Posts.SetCurrentTab(endPoint.String())
	d.App.Draw(d.Posts)
}

func (d *Display) generatePostsList(point api.EndPoint) *cview.List {

	posts := cview.NewList()
	posts.SetPadding(0, 1, 0, 0)
	posts.SetHighlightFullLine(false)
	posts.SetChangedFunc(d.postsHandler)

	collectionIDs, err := d.DB.Collection(point)
	if err != nil {
		log.Fatal(err)
	}

	d.StatusBar.SetMax(len(collectionIDs) - 1)

	subscription := d.DB.Subscribe(collectionIDs)

	renderlist := func() {
		var shortcutCounter = 0
		for post := range subscription.Updates() {
			item := cview.NewListItem(formatPrimaryLine(post, false))
			item.SetReference(post)
			item.SetShortcut(rune('a' + shortcutCounter))
			shortcutCounter++
			item.SetSelectedFunc(func() {
				d.App.SetFocus(d.Comments.Tree)
			})
			posts.AddItem(item)

			if d.StatusBar.Complete() {
				d.StatusBar.SetProgress(0)
			} else {
				d.StatusBar.AddProgress(1)
			}
			d.App.QueueUpdateDraw(func() {})
		}
	}

	go renderlist()

	return posts
}

func (d *Display) postsHandler(_ int, listItem *cview.ListItem) {
	post := listItem.GetReference().(api.Post)
	commentsRoot := d.generateCommentTree(post)
	d.Comments.Tree.SetRoot(commentsRoot)
	d.Comments.Tree.SetTopLevel(0)
}

func formatPrimaryLine(post api.Post, selected bool) string {
	link, err := url.Parse(post.URL())
	if err != nil {
		log.Print(err)
	}

	var points string

	if post.Score() > 1 {
		points = fmt.Sprintf("[#8a8a8a:-:d]%d points", post.Score())
	} else {
		points = fmt.Sprintf("[#555555:-:d]%d point", post.Score())
	}

	if selected {
		return fmt.Sprintf("%s [#8a8a8a:-:b](%s) [-:-:]by %s -- %s", post.Title(), link.Host, post.By(), points)
	} else {
		return fmt.Sprintf("%s [#8a8a8a:-:d](%s) [-:-:]by %s -- %s", post.Title(), link.Host, post.By(), points)
	}

}

//func formatSecondaryLine(post api.Post) string {
//	var kids, points string
//	if len(post.Kids()) > 1 {
//		kids = fmt.Sprintf("[#8a8a8a:-:b]%d comments", len(post.Kids()))
//	} else {
//		kids = fmt.Sprintf("[#555555:-:d]%d comment", len(post.Kids()))
//	}
//	if post.Score() > 1 {
//		points = fmt.Sprintf("[#8a8a8a:-:d] %d points", post.Score())
//	} else {
//		points = fmt.Sprintf("[#555555:-:d] %d point", post.Score())
//	}
//
//
//
//
//	return fmt.Sprintf("%s [-:-:d]by [#ff6600:-:d]%s [-:-:d]%s ago [-:-:-]| %s",  points, post.By(), humanize.Time(post.Time()), kids)
//
//}
