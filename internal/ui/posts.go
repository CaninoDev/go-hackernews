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
	ids   map[string][]int
	lists map[string]*cview.List
}

func (d *Display) NewPostsView() *Posts {
	tabbedPanels := cview.NewTabbedPanels()
	postIDs := make(map[string][]int)
	lists := make(map[string]*cview.List)


	for _, endPoint := range d.DB.EndPoints() {
		lists[endPoint] = cview.NewList()
		lists[endPoint].SetPadding(0, 1, 0, 0)
		lists[endPoint].SetHighlightFullLine(false)
		listItem := cview.NewListItem("")
		lists[endPoint].AddItem(listItem)
		go d.generatePostIDs(endPoint)
		tabbedPanels.AddTab(endPoint, endPoint, lists[endPoint])
	}

	//tabbedPanels.SetChangedFunc(d.generatePostList)
	posts := &Posts{
		tabbedPanels,
		postIDs,
		lists,
	}
	return posts
}

func (d *Display) generatePostIDs(endpoint string) {
	d.Posts.ids[endpoint], _ = d.DB.Collection(endpoint)
}

func (d *Display) generatePostList() {
	subscription := d.DB.Subscribe(d.Posts.ids[d.Posts.GetCurrentTab()])
	d.StatusBar.SetMax(len(d.Posts.ids[d.Posts.GetCurrentTab()]))
	var shortcutCounter = 0
	renderList := func() {
		for post := range subscription.Updates() {
			item := cview.NewListItem(formatPrimaryLine(post, false))
			item.SetReference(post)
			item.SetShortcut(rune('a' + shortcutCounter))
			shortcutCounter++
			item.SetSelectedFunc(func() {
				d.App.SetFocus(d.Comments.Tree)
			})
			d.Posts.lists[d.Posts.GetCurrentTab()].AddItem(item)
			if d.StatusBar.Complete() {
				d.StatusBar.SetProgress(0)
			} else {
				d.StatusBar.AddProgress(1)
			}
			d.App.QueueUpdateDraw(func() {})
		}
	}
	go renderList()

	d.App.SetFocus(d.Posts)
	//d.Posts.ids[currentList] = d.DB.Collection()
	//subscription := d.DB.Subscribe(d.Posts.ids[currentList])
	//d.StatusBar.SetMax(len(d.Posts.ids[currentList]) - 1)


	//d.DebugBar.SetText(fmt.Sprintf("%d",len(d.Posts.ids[currentList]) - 1))
	//var shortcutCounter = 0
	//render := func() {
	//	for post := range subscription.Updates() {
	//		d.DebugBar.Clear()
	//		//item := cview.NewListItem(formatPrimaryLine(post, false))
	//		//item.SetReference(post)
	//		//item.SetShortcut(rune('a' + shortcutCounter))
	//		//shortcutCounter++
	//		//item.SetSelectedFunc(func() {
	//		//	d.App.SetFocus(d.Comments.Tree)
	//		//})
	//		//d.Posts.lists[currentList].AddItem(item)
	//		priorText := d.DebugBar.GetText(false)
	//		d.DebugBar.SetText(fmt.Sprintf("%s, %s", priorText, post.Title()))
	//		if d.StatusBar.Complete() {
	//			d.StatusBar.SetProgress(0)
	//		} else {
	//			d.StatusBar.AddProgress(1)
	//		}
	//
	//	}
	//}
	//go d.App.QueueUpdateDraw(render)
}
//
//func (d *Display) generatePostList() {
//	currentEndPoint := d.Posts.GetCurrentTab()
//
//	d.StatusBar.SetMax(len(d.Posts.ids[currentEndPoint]) - 1)
//	subscription := d.DB.Subscribe(d.Posts.ids[currentEndPoint])
//
//	renderlist := func() {
//		var shortcutCounter = 0
//		for post := range subscription.Updates() {
//			item := cview.NewListItem(formatPrimaryLine(post, false))
//			item.SetReference(post)
//			item.SetShortcut(rune('a' + shortcutCounter))
//			shortcutCounter++
//			item.SetSelectedFunc(func() {
//				d.App.SetFocus(d.Comments.Tree)
//			})
//			d.Posts.lists[currentEndPoint].AddItem(item)
//
//			if d.StatusBar.Complete() {
//				d.StatusBar.SetProgress(0)
//			} else {
//				d.StatusBar.AddProgress(1)
//			}
//			d.App.QueueUpdateDraw(func() {})
//		}
//	}
//
//	go renderlist()
//}
func (d *Display) generatePostsIDs() map[string][]int {
	posts := make(map[string][]int)
	endPoints := []api.EndPoint{
		api.NewS, api.Top, api.Best, api.Ask, api.Show,
	}

	for _, endPoint := range endPoints {
		posts[endPoint.String()] = []int{0}
	}
	return posts
}

//func (d *Display) NewPosts() {
//	d.Posts = make(map[string]*cview.List)
//	endPoints := []api.EndPoint{
//		api.NewS, api.Top, api.Best, api.Ask, api.Show,
//	}
//	populateList := func() {
//		for _, endPoint := range endPoints {
//			d.generatePostsList(endPoint)
//		}
//	}
//	go populateList()
//}

func (d *Display) List(endPoint api.EndPoint) {
	d.Posts.SetCurrentTab(endPoint.String())
}

//
//func (d *Display) generatePostsList(point api.EndPoint) {
//	d.Posts[point.String()].SetPadding(0, 1, 0, 0)
//	d.Posts[point.String()].SetHighlightFullLine(false)
//	d.Posts[point.String()].SetChangedFunc(d.postsHandler)
//
//	collectionIDs, err := d.DB.Collection(point)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	d.StatusBar.SetMax(len(collectionIDs) - 1)
//
//	subscription := d.DB.Subscribe(collectionIDs)
//
//	renderlist := func() {
//		var shortcutCounter = 0
//		for post := range subscription.Updates() {
//			item := cview.NewListItem(formatPrimaryLine(post, false))
//			item.SetReference(post)
//			item.SetShortcut(rune('a' + shortcutCounter))
//			shortcutCounter++
//			item.SetSelectedFunc(func() {
//				d.App.SetFocus(d.Comments.Tree)
//			})
//			d.Posts[point.String()].AddItem(item)
//
//			if d.StatusBar.Complete() {
//				d.StatusBar.SetProgress(0)
//			} else {
//				d.StatusBar.AddProgress(1)
//			}
//			d.App.QueueUpdateDraw(func() {})
//		}
//	}
//
//	go renderlist()
//
//	d.Posts.AddTab(point.String(), point.String(), d.Posts[point.String()])
//}

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
