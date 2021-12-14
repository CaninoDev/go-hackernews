package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"fmt"
	"github.com/CaninoDev/go-hackernews/internal/api"
	"github.com/CaninoDev/go-hackernews/internal/store"
	"log"
	"net/url"
)

type Browser struct {
	app        *App
	lists      *cview.TabbedPanels
	states     map[string]ListState
	currentTab string
	statusBar  *StatusBar
	debugBar   *cview.TextView
}

type ListState struct {
	*cview.List
	ids              []int
	currentPageIndex int
}

func NewBrowser(app *App) *Browser {
	tabbedList := cview.NewTabbedPanels()
	return &Browser{
		app:        app,
		lists:      tabbedList,
		states:     make(map[string]ListState),
		currentTab: "",
		statusBar:  NewStatusBar(),
		debugBar:   NewDebugBar(),
	}
}

func (b *Browser) initializeTabbedLists() {
	tabLabels := b.app.store.CollectionsList()
	b.states = make(map[string]ListState)
	for _, tabLabel := range tabLabels {
		if b.currentTab == "" {
			b.currentTab = tabLabel
		}
		list := cview.NewList()
		list.SetPadding(1, 1, 0, 0)
		list.SetHighlightFullLine(false)
		list.SetSelectedFunc(b.listItemHandler)
		b.states[tabLabel] = ListState{
			List:             list,
			ids:              b.app.store.Collection(tabLabel),
			currentPageIndex: 0,
		}
		b.lists.AddTab(tabLabel, tabLabel, b.states[tabLabel])
	}
	//b.lists.SetChangedFunc(b.populateList)
}

func (b *Browser) listItemHandler(_ int, listItem *cview.ListItem) {
	post := listItem.GetReference().(store.Item)
	b.debugBar.SetText(post.Title())
}

func (b *Browser) populateList() {
	currentTab := b.lists.GetCurrentTab()

	pagedBatch := b.paginate(currentTab)

	b.statusBar.SetMax(len(pagedBatch) - 1)
	for _, id := range pagedBatch {
		post, err := b.app.store.Item(id)
		if err != nil {
			b.debugBar.SetText(fmt.Sprintf("%v", err))
		}
		listItem := cview.NewListItem(formatPrimaryLine(post, false))
		listItem.SetReference(post)
		b.states[currentTab].AddItem(listItem)
		b.app.ui.QueueUpdateDraw(func() {})
		b.statusBar.AddProgress(1)
		b.app.ui.QueueUpdateDraw(func() {})
	}
	if b.statusBar.Complete() {
		b.statusBar.SetProgress(0)
	}

	//go func() {
	//	for id := range pagedBatch {
	//		post, err := b.app.store.Item(id)
	//		if err != nil {
	//			b.debugBar.SetText(fmt.Sprintf("%v", err))
	//		}
	//		listItem := cview.NewListItem(formatPrimaryLine(post, false))
	//		listItem.SetReference(post)
	//		b.states[currentTab].AddItem(listItem)
	//		b.app.ui.QueueUpdateDraw(func() {})
	//		b.statusBar.AddProgress(1)
	//		b.app.ui.QueueUpdateDraw(func() {})
	//
	//	}
	//	if b.statusBar.Complete() {
	//		b.statusBar.SetProgress(0)
	//		b.app.ui.QueueUpdateDraw(func() {})
	//
	//	}
	//}()
}

func (b *Browser) paginate(currentTab string) []int {

	screenHeight, _ := b.app.ui.GetScreenSize()

	listLength := (screenHeight / 2) - 23
	b.debugBar.SetText(fmt.Sprintf("screenheight: %d, listlength: %d", screenHeight, listLength))
	totalPostCount := len(b.states[currentTab].ids)

	startIndex := (totalPostCount / listLength) * b.states[currentTab].currentPageIndex
	var lastIndex int
	if (startIndex + listLength) > len(b.states[currentTab].ids) {
		lastIndex = len(b.states[currentTab].ids) - 1
	} else {
		lastIndex = startIndex + listLength
	}
	listBatch := b.states[currentTab].ids[startIndex:(lastIndex)]
	return listBatch
}

//func (b *Browser) populate() {
//b.currentTab = b.lists.GetCurrentTab()
//listStr := b.currentTab
//
//listRectHeight, _ := b.app.ui.GetScreenSize()
//
//listLength := (((listRectHeight - 8) / 2) - 4)
//postCount := len(b.states[listStr].ids)
//
//listIndex := (postCount / listLength) * b.states[listStr].currentPageIndex
//
////b.debugBar.SetText(fmt.Sprintf("%d", listIndex))
//totalPostCount := len(b.states[listStr].ids)
//maxIndex := int(math.Ceil(float64(totalPostCount) / float64(listLength)))
//postIDBatch := b.states[listStr].ids[listIndex:(listIndex + listLength)]
//
//b.debugBar.SetText(fmt.Sprintf("#items per page: %d, total #items: %d, currentPage: %d, totalPages:%d", listLength, totalPostCount, b.states[listStr].currentPageIndex, maxIndex))
//b.statusBar.SetMax(len(postIDBatch) - 1)
//for id := range postIDBatch {
//	post, err := b.app.store.Item(id)
//	if err != nil {
//		b.debugBar.SetText(fmt.Sprint(err))
//	}
//	listItem := cview.NewListItem(formatPrimaryLine(post, false))
//	listItem.SetReference(post)
//	b.states[listStr].AddItem(listItem)
//	b.statusBar.AddProgress(1)
//}
//b.statusBar.SetProgress(0)
//}

func (b *Browser) pageNav(next bool) {
	//listStr := b.currentTab
	//oldCurrentIndex := b.states[b.currentTab].currentPageIndex
	//
	//listRectHeight, _ := b.app.ui.GetScreenSize()
	//
	//listLength := (listRectHeight - 12) / 2
	//totalPostCount := len(b.states[listStr].ids)
	//maxIndex := int(math.Ceil(float64(totalPostCount) / float64(listLength)))
	//
	//if next {
	//	if oldCurrentIndex >= maxIndex {
	//		lastItemIndex := b.states[listStr].GetItemCount() - 1
	//		b.states[listStr].SetCurrentItem(lastItemIndex)
	//	} else {
	//		//b.states[listStr].currentPageIndex++
	//		b.states[listStr].Clear()
	//		b.app.ui.QueueUpdateDraw(b.populate, b.lists)
	//	}
	//} else {
	//	if oldCurrentIndex == 1 {
	//		b.states[listStr].SetCurrentItem(1)
	//	} else {
	//		//b.states[listStr].currentPageIndex--
	//		b.states[listStr].Clear()
	//		b.app.ui.QueueUpdateDraw(b.populate, b.lists)
	//	}
	//}

	//b.debugBar.SetText(fmt.Sprintf("Items[#CurrentList:%d -- #TotalItems:%d] Pages[Old:%d -- New:%d -- TotalPages:%d]", listLength, len(b.states[listStr].ids), oldCurrentIndex, b.states[listStr].currentPageIndex, maxIndex))
}

//subscription := b.app.store.Subscribe(postIDBatch)
//for post := range subscription.Updates() {
//	listItem := cview.NewListItem(formatPrimaryLine(post, false))
//	listItem.SetShortcut(rune('a' + shortCutCounter))
//	shortCutCounter++
//	listItem.SetReference(post)
//	b.states[listStr].AddItem(listItem)
//	if b.statusBar.Complete() {
//		b.statusBar.SetProgress(0)
//
//	} else {
//		b.statusBar.AddProgress(1)
//	}
//	b.app.ui.QueueUpdateDraw(func() {})
//}
//listLength := listRectHeight/2 - 1
//postCount := len(b.states[listStr].listItemsCache)
//
//listIndex := (postCount / listLength) * b.states[listStr].currentPageIndex
//postIDBatch := b.states[listStr].listItemsCache[listIndex:(listIndex + listLength)]
//for _, id := range postIDBatch {
//	addListItem := func() {
//		post, err := b.app.store.Item(id)
//		if err != nil {
//			b.debugWrite(fmt.Sprintf("%v", err))
//		}
//		listItem := cview.NewListItem(formatPrimaryLine(post, false))
//		listItem.SetReference(post)
//		b.states[listStr].AddItem(listItem)
//		if b.statusBar.Complete() {
//			b.statusBar.SetProgress(0)
//		} else {
//			b.statusBar.AddProgress(1)
//		}
//	}
//	b.app.ui.QueueUpdate(addListItem)
//}

//func (a *App) generateList() {
//	endPointStr := a.ListPanels.GetCurrentTab()
//	endPoint := api.ToEndPoint(endPointStr)
//	ids := a.store.Collection(endPoint)
//	var shortcutCounter = 0
//	a.StatusBar.SetMax(len(ids))
//	subscription := a.store.Subscribe(ids)
//	for post := range subscription.Updates() {
//		listItem := cview.NewListItem(formatPrimaryLine(post, false))
//		listItem.SetReference(post)
//		listItem.SetShortcut(rune('a' + shortcutCounter))
//		shortcutCounter++
//		a.ListPanels.posts[endPoint].AddItem(listItem)
//		if a.StatusBar.Complete() {
//			a.StatusBar.SetProgress(0)
//		} else {
//			a.StatusBar.AddProgress(1)
//		}
//	}
//}

//
//func (d *ui) generatePosts(endpoint string) {
//	var shortcutCounter = 0
//	subscription := d.store.Subscribe(d.store.CollectionIDs(endpoint))
//	for post := range subscription.Updates() {
//		item := cview.NewListItem(formatPrimaryLine(post, false))
//		item.SetReference(post)
//		item.SetShortcut(rune('a' + shortcutCounter))
//		shortcutCounter++
//		d.ListPanels.lists[d.ListPanels.GetCurrentTab()].AddItem(item)
//		if d.StatusBar.Complete() {
//			d.StatusBar.SetProgress(0)
//		} else {
//			d.StatusBar.AddProgress(1)
//		}
//	}
//}
//
//func (d *ui) generatePostsList() {
//
//		subscription := d.store.Subscribe(d.ListPanels.ids[d.ListPanels.GetCurrentTab()])
//		d.DebugBar.SetText(fmt.Sprintln(len(d.ListPanels.ids[d.ListPanels.GetCurrentTab()])))
//		var shortcutCounter = 0
//		for post := range subscription.Updates() {
//			item := cview.NewListItem(formatPrimaryLine(post, false))
//			item.SetReference(post)
//			item.SetShortcut(rune('a' + shortcutCounter))
//			shortcutCounter++
//			d.ListPanels.lists[d.ListPanels.GetCurrentTab()].AddItem(item)
//			if d.StatusBar.Complete() {
//				d.StatusBar.SetProgress(0)
//			} else {
//				d.StatusBar.AddProgress(1)
//			}
//			d.ui.QueueUpdateDraw(func() {})
//		}
//	d.ui.SetFocus(d.ListPanels)
//}
//
//func (d *ui) generatePostList() {
//	renderList := func() {
//		subscription := d.store.Subscribe(d.ListPanels.ids[d.ListPanels.GetCurrentTab()])
//		d.StatusBar.SetMax(len(d.ListPanels.ids[d.ListPanels.GetCurrentTab()]))
//		var shortcutCounter = 0
//		for post := range subscription.Updates() {
//			item := cview.NewListItem(formatPrimaryLine(post, false))
//			item.SetReference(post)
//			item.SetShortcut(rune('a' + shortcutCounter))
//			shortcutCounter++
//			item.SetSelectedFunc(func() {
//				d.ui.SetFocus(d.CommentsTree.Tree)
//			})
//			d.ListPanels.lists[d.ListPanels.GetCurrentTab()].AddItem(item)
//			if d.StatusBar.Complete() {
//				d.StatusBar.SetProgress(0)
//			} else {
//				d.StatusBar.AddProgress(1)
//			}
//			d.ui.QueueUpdateDraw(func() {})
//		}
//	}
//	go renderList()
//	d.ui.SetFocus(d.ListPanels)
//}
//
//func (d *ui) generatePostList() {
//	currentEndPoint := d.ListPanels.GetCurrentTab()
//
//	d.StatusBar.SetMax(len(d.ListPanels.ids[currentEndPoint]) - 1)
//	subscription := d.store.Subscribe(d.ListPanels.ids[currentEndPoint])
//
//	renderlist := func() {
//		var shortcutCounter = 0
//		for post := range subscription.Updates() {
//			item := cview.NewListItem(formatPrimaryLine(post, false))
//			item.SetReference(post)	wg.Wait()
//			item.SetShortcut(rune('a' + shortcutCounter))
//			shortcutCounter++
//			item.SetSelectedFunc(func() {
//				d.ui.SetFocus(d.CommentsTree.Tree)
//			})
//			d.ListPanels.lists[currentEndPoint].AddItem(item)
//
//			if d.StatusBar.Complete() {
//				d.StatusBar.SetProgress(0)
//			} else {
//				d.StatusBar.AddProgress(1)
//			}
//			d.ui.QueueUpdateDraw(func() {})
//		}
//	}
//
//	go renderlist()

//func (d *ui) generatePostsList(point api.EndPoint) {
//	d.ListPanels.posts[point.String()].SetPadding(0, 1, 0, 0)
//	d.ListPanels.posts[point.String()].SetHighlightFullLine(false)
//	d.ListPanels.posts[point.String()].SetChangedFunc(d.postsHandler)
//
//	collectionIDs, err := d.store.CollectionIDs(point)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	d.StatusBar.SetMax(len(collectionIDs) - 1)
//
//	subscription := d.store.Subscribe(collectionIDs)
//
//	renderlist := func() {
//		var shortcutCounter = 0
//		for post := range subscription.Updates() {
//			item := cview.NewListItem(formatPrimaryLine(post, false))
//			item.SetReference(post)
//			item.SetShortcut(rune('a' + shortcutCounter))
//			shortcutCounter++
//			item.SetSelectedFunc(func() {
//				d.ui.SetFocus(d.CommentsTree.Tree)
//			})
//			d.ListPanels[point.String()].AddItem(item)
//
//			if d.StatusBar.Complete() {
//				d.StatusBar.SetProgress(0)
//			} else {
//				d.StatusBar.AddProgress(1)
//			}
//			d.ui.QueueUpdateDraw(func() {})
//		}
//	}
//
//	go renderlist()
//
//	d.ListPanels.AddTab(point.String(), point.String(), d.ListPanels[point.String()])
//}

//func (a *App) postsHandler(_ int, listItem *cview.ListItem) {
//	post := listItem.GetReference().(store.Item)
//	commentsRoot := a.generateCommentTree(post)
//	a.CommentsTree.Tree.SetRoot(commentsRoot)
//	a.CommentsTree.Tree.SetTopLevel(0)
//}

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
