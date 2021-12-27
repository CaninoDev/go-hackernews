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

type Nav int

const (
	prev Nav = iota
	next
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
	maxListItemWidth      int
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
	_, _, w, _ := t.tabbedLists.GetInnerRect()
	setReadTimeStamp := func() {
		t.app.store.SetItemReadStamp(&post)

		listItem.SetReference(post)
		fmtStr, _ := formatReadPostLine(post, w)
		listItem.SetMainText(fmtStr)
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

		width := t.app.width

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

			fmtStr, maxLenStr := formatPrimaryLine(post, width)

			listItem := cview.NewListItem(fmtStr)
			if t.states[currentTab].maxListItemWidth < maxLenStr {
				t.states[currentTab].maxListItemWidth = maxLenStr
			}

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

func (t *TabbedLists) resizeListItems(width int) {
	currentTab := t.tabbedLists.GetCurrentTab()
	if width < t.states[currentTab].maxListItemWidth {
		listItemCount := t.states[currentTab].GetItemCount()
		for i := 0; i < listItemCount; i++ {
			listItem := t.states[currentTab].GetItem(i)
			postRef := listItem.GetReference()
			post := postRef.(store.Item)
			fmtStr, _ := formatPrimaryLine(post, width)
			resizeListItemText := func() {
				listItem.SetMainText(fmtStr)
			}
			t.app.ui.QueueUpdateDraw(resizeListItemText)
		}
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

// pageNav will trigger the visible list to the nav page of available items (true)
// or the prior page based on the provided flag.
func (t *TabbedLists) pageNav(nav Nav) {
	currentTab := t.tabbedLists.GetCurrentTab()

	// Capture the current state of the list
	currentPageIndex := &t.states[currentTab].currentPageIndex
	listLength := t.states[currentTab].GetItemCount()
	totalPostCount := len(t.states[currentTab].itemIDs)
	maxIndex := int(math.Ceil(float64(totalPostCount) / float64(listLength)))

	switch nav {
	case next:
		// If we are on the last page, select the last item on the list. Otherwise
		// repopulate the list with the contents of the next page.
		if *currentPageIndex >= maxIndex {
			lastItemIndex := t.states[currentTab].GetItemCount()
			t.states[currentTab].SetCurrentItem(lastItemIndex)
		} else {
			*currentPageIndex++
			t.populateList()
		}
	case prev:
		// If we are on the first page, select the first item on the list. Otherwise
		// repopulate list with the contents of the previous page.
		if *currentPageIndex == 0 {
			t.states[currentTab].SetCurrentItem(0)
		} else {
			*currentPageIndex--
			t.populateList()
		}
	}

	//t.app.statusBar.SetText(
	//	fmt.Sprintf(
	//		"Items[#CurrentList:%d -- #TotalItems:%d] Pages[Old:%d -- New:%d -- TotalPages:%d]",
	//		listLength,
	//		len(t.states[currentTab].itemIDs),
	//		*currentPageIndex,
	//		t.states[currentTab].currentPageIndex,
	//		maxIndex,
	//	),
	//)
}

// formatPrimaryLine will return a formmatted string for the item's title
// based on various attributes.
func formatPrimaryLine(post store.Item, terminal_width int) (string, int) {
	if post.GetReadStamp() != time.Unix(0, 0) {
		return formatReadPostLine(post, terminal_width)
	} else {
		return formatUnreadPostLine(post, terminal_width)
	}
}

// formatReadPostLine will format the postView's title string if tagged as read.
func formatReadPostLine(post store.Item, tWidth int) (string, int) {
	link, err := url.Parse(post.URL())
	if err != nil {
		log.Print(err)
	}

	var pLen int
	points := formatPoints(post.Score())

	if post.Score() <= 1 {
		pLen = 7
	} else if pLen <= 9 {
		pLen = 8
	} else {
		pLen = 9
	}

	pBy := len(post.By())
	pTitle := len(post.Title())

	if len(link.Host) > 0 {
		pHost := len(link.Host)
		eWidth := 11 + pHost + pBy + pLen
		stringLength := eWidth + pTitle
		if stringLength > tWidth {
			oSpaces := stringLength - tWidth
			if oSpaces < 3 {
				oSpaces = 3
			}
			oldTitle := post.Title()
			truncatedTitle := oldTitle[:(pTitle-oSpaces)] + "..."
			return fmt.Sprintf("[-:-:d]%s (%s) [::-]by %s -- %s", truncatedTitle, link.Host, post.By(), points), stringLength
		} else {
			return fmt.Sprintf("[-:-:d]%s (%s) [::-]by %s -- %s", post.Title(), link.Host, post.By(), points), stringLength
		}
	} else {
		eWidth := 8 + pBy + pLen
		stringLength := eWidth + pTitle
		if stringLength > tWidth {
			oSpaces := tWidth - eWidth
			if oSpaces < 3 {
				oSpaces = 3
			}
			oldTitle := post.Title()
			truncatedTitle := oldTitle[:oSpaces] + "..."
			return fmt.Sprintf("[-:-:d]%s [::-]by %s -- %s", truncatedTitle, post.By(), points), stringLength
		}
		return fmt.Sprintf("[-:-:d]%s [::-]by %s -- %s", post.Title(), post.By(), points), stringLength
	}
}

// formatReadPostLine will format the postView's title string if items was unread.
func formatUnreadPostLine(post store.Item, tWidth int) (string, int) {
	link, err := url.Parse(post.URL())
	if err != nil {
		log.Print(err)
	}

	var pLen int
	points := formatPoints(post.Score())

	if post.Score() <= 1 {
		pLen = 7
	} else if post.Score() <= 9 {
		pLen = 8
	} else {
		pLen = 9
	}

	pBy := len(post.By())
	pTitle := len(post.Title())

	if len(link.Host) > 0 {
		pHost := len(link.Host)
		eWidth := pHost + pBy + pLen + 16
		stringLength := eWidth + pTitle
		if stringLength > tWidth {
			oSpaces := stringLength - tWidth
			if oSpaces < 3 {
				oSpaces = 3
			}
			oldTitle := post.Title()
			title := oldTitle[:(pTitle-oSpaces)] + "..."
			return fmt.Sprintf("[-:-:b]%s (%s) [::-]by %s -- %s", title, link.Host, post.By(), points), stringLength
		} else {
			return fmt.Sprintf("[-:-:b]%s (%s) [::-]by %s -- %s", post.Title(), link.Host, post.By(), points), stringLength
		}
	} else {
		eWidth := 8 + pBy + pLen
		stringLength := eWidth + pTitle
		if (pTitle + eWidth) > tWidth {
			oSpaces := tWidth - eWidth - 3
			title := post.Title()[:oSpaces] + "..."
			return fmt.Sprintf("[-:-:d]%s [::-]by %s -- %s", title, post.By(), points), stringLength
		}
		return fmt.Sprintf("[-:-:d]%s [::-]by %s -- %s", post.Title(), post.By(), points), stringLength
	}
}

// formatPoints will return a color-graded string based on the provided score.
func formatPoints(score int) string {
	colorScale := []string{
		"#11174B", "#162065", "#1C2B7F", "#24448E", "#2D5E9E", "#3577AE", "#3D91BE", "#46ACE0", "#62BED2", "#8ACDCE", "#B3DDCC", "#DCECC9",
	}
	switch {
	case score <= 1:
		return fmt.Sprintf("[%s::]%d point", colorScale[0], score)
	case score <= 7:
		return fmt.Sprintf("[%s::]%d points", colorScale[0], score)
	case score <= 15:
		return fmt.Sprintf("[%s::]%d points", colorScale[1], score)
	case score <= 23:
		return fmt.Sprintf("[%s::]%d points", colorScale[2], score)
	case score <= 31:
		return fmt.Sprintf("[%s::]%d points", colorScale[3], score)
	case score <= 39:
		return fmt.Sprintf("[%s::]%d points", colorScale[4], score)
	case score <= 47:
		return fmt.Sprintf("[%s::]%d points", colorScale[5], score)
	case score <= 55:
		return fmt.Sprintf("[%s::]%d points", colorScale[6], score)
	case score <= 63:
		return fmt.Sprintf("[%s::]%d points", colorScale[7], score)
	case score <= 71:
		return fmt.Sprintf("[%s::]%d points", colorScale[8], score)
	case score <= 79:
		return fmt.Sprintf("[%s::]%d points", colorScale[9], score)
	case score <= 87:
		return fmt.Sprintf("[%s::]%d points", colorScale[10], score)
	default:
		return fmt.Sprintf("[%s::]%d points", colorScale[11], score)
	}
}
