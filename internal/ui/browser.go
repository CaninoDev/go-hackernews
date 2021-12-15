package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"fmt"
	"github.com/CaninoDev/go-hackernews/internal/api"
	"github.com/CaninoDev/go-hackernews/internal/store"
	"log"
	"math"
	"net/url"
)

type Browser struct {
	app        *App
	lists      *cview.TabbedPanels
	states     map[string]*listState
	currentTab string
	statusBar  *StatusBar
	debugBar   *cview.TextView
}

type listState struct {
	*cview.List
	ids              []int
	currentPageIndex int
}

func NewBrowser(app *App) *Browser {
	tabbedList := cview.NewTabbedPanels()
	return &Browser{
		app:        app,
		lists:      tabbedList,
		states:     make(map[string]*listState),
		currentTab: "",
		statusBar:  NewStatusBar(),
		debugBar:   NewDebugBar(),
	}
}

func (b *Browser) initializeTabbedLists() {
	tabLabels := b.app.store.CollectionsList()
	b.states = make(map[string]*listState)
	for _, tabLabel := range tabLabels {
		if b.currentTab == "" {
			b.currentTab = tabLabel
		}
		list := cview.NewList()
		list.SetPadding(1, 1, 0, 0)
		list.SetHighlightFullLine(false)
		list.SetSelectedFunc(b.listItemHandler)
		b.states[tabLabel] = &listState{
			List:             list,
			ids:              b.app.store.Collection(tabLabel),
			currentPageIndex: 0,
		}
		b.lists.AddTab(tabLabel, tabLabel, b.states[tabLabel])
	}

}

func (b *Browser) listItemHandler(_ int, listItem *cview.ListItem) {
	post := listItem.GetReference().(store.Item)
	b.debugBar.SetText(post.Title())
}

func (b *Browser) populateList() {
	currentTab := b.lists.GetCurrentTab()

	pagedBatch, totalPages := b.paginate(currentTab)

	tabLabel := fmt.Sprintf("%s(%d/%d)", currentTab, b.states[currentTab].currentPageIndex, totalPages)
	b.lists.SetTabLabel(currentTab, tabLabel)
	b.statusBar.SetMax(len(pagedBatch) - 1)
	b.states[currentTab].Clear()
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
}

func (b *Browser) paginate(currentTab string) ([]int, int) {

	_, _, _, screenHeight := b.lists.GetInnerRect()

	listLength := (screenHeight / 2) - 2
	b.debugBar.SetText(fmt.Sprintf("screenheight: %d, listlength: %d", screenHeight, listLength))
	totalPostCount := len(b.states[currentTab].ids)

	totalPageCount := math.Ceil(float64(totalPostCount) / float64(listLength))

	startIndex := int(totalPageCount) * b.states[currentTab].currentPageIndex
	var lastIndex int
	if (startIndex + listLength) > len(b.states[currentTab].ids) {
		lastIndex = len(b.states[currentTab].ids) - 1
	} else {
		lastIndex = startIndex + listLength
	}
	listBatch := b.states[currentTab].ids[startIndex:(lastIndex)]
	return listBatch, int(totalPageCount)
}

func (b *Browser) pageNav(next bool) {

	listStr := b.currentTab
	oldCurrentIndex := b.states[b.currentTab].currentPageIndex

	listRectHeight, _ := b.app.ui.GetScreenSize()

	listLength := (listRectHeight - 12) / 2
	totalPostCount := len(b.states[listStr].ids)
	maxIndex := int(math.Ceil(float64(totalPostCount) / float64(listLength)))

	if next {
		if oldCurrentIndex >= maxIndex {
			lastItemIndex := b.states[listStr].GetItemCount() - 1
			b.states[listStr].SetCurrentItem(lastItemIndex)
		} else {
			b.states[listStr].currentPageIndex++
			b.populateList()
		}
	} else {
		if oldCurrentIndex == 1 {
			b.states[listStr].SetCurrentItem(1)
		} else {
			b.states[listStr].currentPageIndex--
			b.populateList()
		}
	}

	b.debugBar.SetText(fmt.Sprintf("Items[#CurrentList:%d -- #TotalItems:%d] Pages[Old:%d -- New:%d -- TotalPages:%d]", listLength, len(b.states[listStr].ids), oldCurrentIndex, b.states[listStr].currentPageIndex, maxIndex))
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
