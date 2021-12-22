package ui

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"sync"
	"time"

	"code.rocketnine.space/tslocum/cview"
	"github.com/CaninoDev/go-hackernews/internal/store"
)

// TabbedList contains the various primitives and states necessary to
// provide a set of tabbed lists of posts.
type TabbedLists struct {
	app         *App
	tabbedLists *cview.TabbedPanels
	states      map[string]*listState
	statusBar   *cview.ProgressBar
	sync.RWMutex
}

// lisState provides a pointer to a list and state for each
// tab in the TabbedLists.
type listState struct {
	*cview.List
	itemIDs               []int
	currentPageIndex      int
	lastSelectedItemIndex int
}

// NewTabbedLists returns a TabbedList with an
// embedded pointer to App.
func NewTabbedLists(app *App) *TabbedLists {
	tabbedList := cview.NewTabbedPanels()
	return &TabbedLists{
		app:         app,
		tabbedLists: tabbedList,
		states:      make(map[string]*listState),
		statusBar:   cview.NewProgressBar(),
	}
}

// initializeTabbedLists initializes the TabbedView primitive with the available endpoints
// and defines the various behaviors of the list.
func (t *TabbedLists) initializeTabbedLists() {
	tabLabels := t.app.store.CollectionsList()

	for _, tabLabel := range tabLabels {

		list := cview.NewList()
		list.SetPadding(1, 1, 0, 0)
		list.SetHighlightFullLine(false)
		list.SetSelectedFunc(t.listItemHandler)
		list.SetPadding(1, 1, 1, 1)
		list.SetCompact(true)
		list.SetScrollBarVisibility(cview.ScrollBarNever)
		list.SetSelectedAlwaysCentered(false)
		list.SetChangedFunc(t.listChangedHandler)
		t.states[tabLabel] = &listState{
			List:                  list,
			itemIDs:               t.app.store.Collection(tabLabel),
			currentPageIndex:      0,
			lastSelectedItemIndex: 0,
		}

		t.tabbedLists.AddTab(tabLabel, tabLabel, t.states[tabLabel])
	}

	t.tabbedLists.SetTabSwitcherAfterContent(true)
	// t.tabbedLists.SetChangedFunc(t.tabsHandler)
}

func (t *TabbedLists) listChangedHandler(idx int, _ *cview.ListItem) {
	t.setLastSelectedItemIndex(idx)
}

//func (t *TabbedLists) tabsHandler() {
//	t.app.ui.QueueUpdateDraw(func() {
//		currentTab := t.tabbedLists.GetCurrentTab()
//		listItemCount := t.states[currentTab].GetItemCount()
//		lastSelectedIndex := t.states[currentTab].lastSelectedItemIndex
//
//		if listItemCount > 0 && lastSelectedIndex != 0 {
//			t.states[currentTab].SetCurrentItem(lastSelectedIndex)
//		}
//	})
//}

// setSLastSelectedItemIndex records the last selected item index.
func (t *TabbedLists) setLastSelectedItemIndex(idx int) {
	t.Lock()
	defer t.Unlock()

	currentTab := t.tabbedLists.GetCurrentTab()
	t.states[currentTab].lastSelectedItemIndex = idx
}

// listItemHandler defines the behvarior when the user selects an item from the list; namely
// sets the item read timestamp, formats the list item accordingly to indicate that it has been
// read, loads the comment tree and viewer, and switches over to the postView view panel.
func (t *TabbedLists) listItemHandler(selectedItemIndex int, listItem *cview.ListItem) {
	post := listItem.GetReference().(store.Item)

	setReadTimeStamp := func() {

		t.app.store.SetItemReadStamp(&post)

		listItem.SetReference(post)

		listItem.SetMainText(formatReadPostLine(post))

	}

	t.app.ui.QueueUpdateDraw(setReadTimeStamp)
	t.app.postView.SetPost(post)
	t.app.panels.HidePanel(LISTPANEL)
	t.app.panels.ShowPanel(POSTPANEL)
}

// populateList determines the currently selected tabbed list, the dimensions of the list's primitive,
// polls the store to populate the list to just fill the available space, and update the progress bar.
func (t *TabbedLists) populateList() {
	currentTab := t.tabbedLists.GetCurrentTab()

	pagedBatch, totalPages := t.paginate(currentTab)

	renderList := func() {
		paginationInfo := fmt.Sprintf("(%d/%d)", t.states[currentTab].currentPageIndex+1, totalPages)

		clearList := func() {
			t.app.statusBar.SetText(paginationInfo)
			t.statusBar.SetMax(len(pagedBatch) - 1)
			t.states[currentTab].lastSelectedItemIndex = 0
			t.states[currentTab].Clear()
		}

		t.app.ui.QueueUpdateDraw(clearList)

		resetProgressBar := func() {
			t.statusBar.SetProgress(0)
		}

		for _, id := range pagedBatch {
			post, err := t.app.store.Item(id)
			if err != nil {
				t.app.statusBar.SetText(fmt.Sprintf("%v", err))
			}
			listItem := cview.NewListItem(formatPrimaryLine(post))
			listItem.SetReference(post)

			// addItem adds the item to the current list and updates the progress bar.
			addItem := func() {
				t.states[currentTab].AddItem(listItem)
				t.statusBar.AddProgress(1)
			}
			t.app.ui.QueueUpdateDraw(addItem)
		}

		if t.statusBar.Complete() {

			t.app.ui.QueueUpdateDraw(resetProgressBar)
		}
	}

	// If the the list is empty, populate it
	if t.states[currentTab].GetItemCount() == 0 {
		renderList()
	}
}

// paginate calculates the length of the list that can be displayed on the screen,
// the batch of ids from the current state of the list index,
// and the total number of screens(pages) it will take to render the entire list.
func (t *TabbedLists) paginate(currentTab string) ([]int, int) {
	_, _, _, rectHeight := t.tabbedLists.GetInnerRect()

	listLength := rectHeight - 3
	// t.app.statusBar.SetText(fmt.Sprintf("screenheight: %d, listlength: %d", rectHeight, listLength))
	totalPostCount := len(t.states[currentTab].itemIDs)

	totalPageCount := math.Ceil(float64(totalPostCount) / float64(listLength))

	startIndex := int(totalPageCount) * t.states[currentTab].currentPageIndex
	var lastIndex int
	if (startIndex + listLength) > len(t.states[currentTab].itemIDs) {
		lastIndex = len(t.states[currentTab].itemIDs) - 1
	} else {
		lastIndex = startIndex + listLength
	}
	listBatch := t.states[currentTab].itemIDs[startIndex:(lastIndex)]
	return listBatch, int(totalPageCount)
}

// pageNav will trigger the visible list to the next page of available items
// or the prior page based on the provided flag.
func (t *TabbedLists) pageNav(next bool) {
	currentTab := t.tabbedLists.GetCurrentTab()

	oldCurrentIndex := t.states[currentTab].currentPageIndex
	listLength := t.states[currentTab].GetItemCount()
	totalPostCount := len(t.states[currentTab].itemIDs)
	maxIndex := int(math.Ceil(float64(totalPostCount) / float64(listLength)))

	if next {
		if oldCurrentIndex >= maxIndex {
			lastItemIndex := t.states[currentTab].GetItemCount() - 1
			t.states[currentTab].SetCurrentItem(lastItemIndex)
		} else {
			t.states[currentTab].currentPageIndex++
			go t.populateList()
		}
	} else {
		if oldCurrentIndex == 1 {
			t.states[currentTab].SetCurrentItem(1)
		} else {
			t.states[currentTab].currentPageIndex--
			go t.populateList()
		}
	}

	t.app.statusBar.SetText(
		fmt.Sprintf(
			"Items[#CurrentList:%d -- #TotalItems:%d] Pages[Old:%d -- New:%d -- TotalPages:%d]",
			listLength,
			len(t.states[currentTab].itemIDs),
			oldCurrentIndex,
			t.states[currentTab].currentPageIndex,
			maxIndex,
		),
	)
}

// formatPrimaryLine will return a formmatted string for the item's title
// based on various attributes.
func formatPrimaryLine(post store.Item) string {
	if post.GetReadStamp() != time.Unix(0, 0) {
		return formatReadPostLine(post)
	} else {
		return formatUnreadPostLine(post)
	}
}

// formatReadPostLine will format the postView's title string if tagged as read.
func formatReadPostLine(post store.Item) string {
	link, err := url.Parse(post.URL())
	if err != nil {
		log.Print(err)
	}

	points := formatPoints(post.Score())

	if len(link.Host) > 0 {
		return fmt.Sprintf("[-:-:d]%s (%s) [::-]by %s -- %s", post.Title(), link.Host, post.By(), points)
	} else {
		return fmt.Sprintf("[-:-:d]%s [::-]by %s -- %s", post.Title(), post.By(), points)
	}
}

// formatReadPostLine will format the postView's title string if items was unread.
func formatUnreadPostLine(post store.Item) string {
	link, err := url.Parse(post.URL())
	if err != nil {
		log.Print(err)
	}

	points := formatPoints(post.Score())

	if len(link.Host) > 0 {
		return fmt.Sprintf("[::b]%s [::-](%s) by %s -- %s", post.Title(), link.Host, post.By(), points)
	} else {
		return fmt.Sprintf("[::b]%s [::-]by %s -- %s", post.Title(), post.By(), points)
	}
}

// formatPoints will return a color-graded string based on the provided score.
func formatPoints(score int) string {
	one := "#4F742C"
	two := "#729633"
	three := "#8DB13A"
	four := "#B0CE3B"
	five := "#C9D841"
	six := "#D5DC4C"
	seven := "#DCE15D"

	switch {
	case score <= 5:
		return fmt.Sprintf("[%s::]%d points", one, score)
	case score <= 10:
		return fmt.Sprintf("[%s::]%d points", two, score)
	case score <= 15:
		return fmt.Sprintf("[%s::]%d points", three, score)
	case score <= 20:
		return fmt.Sprintf("[%s::]%d points", four, score)
	case score <= 25:
		return fmt.Sprintf("[%s::]%d points", five, score)
	case score <= 30:
		return fmt.Sprintf("[%s::]%d points", six, score)
	default:
		return fmt.Sprintf("[%s::]%d points", seven, score)
	}
}
