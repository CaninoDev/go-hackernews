package ui

import (
	"log"

	"code.rocketnine.space/tslocum/cview"
	"github.com/CaninoDev/go-hackernews/internal/store"
)

var (
	LISTPANEL = "listView"
	POSTPANEL = "postView"
)

type App struct {
	ui           *cview.Application
	focusManager *cview.FocusManager
	store        *store.Store
	statusBar    *cview.TextView
	panels       *cview.Panels
	Cover        *cview.Flex
	listView     *TabbedLists
	postView     *Post
	width        int
	height       int
}

func New() (*App, error) {
	app := &App{}
	app.ui = cview.NewApplication()

	cache, err := store.New()
	if err != nil {
		return &App{}, err
	}

	app.store = cache

	app.listView = NewTabbedLists(app)
	app.postView = NewPostView(app)
	app.statusBar = cview.NewTextView()

	app.ui.SetAfterResizeFunc(app.resizeHandler)
	return app, nil
}

func (a *App) Start() error {
	headerLine := a.initializeHeader()
	a.panels = a.initializePanels()

	rootGrid := cview.NewFlex()
	rootGrid.SetBackgroundColor(cview.Styles.PrimitiveBackgroundColor)
	rootGrid.SetBackgroundTransparent(false)
	rootGrid.SetDirection(cview.FlexRow)
	rootGrid.AddItem(headerLine, 1, 0, false)
	rootGrid.AddItem(a.panels, 0, 1, true)
	a.ui.SetInputCapture(a.inputHandler)
	a.ui.SetRoot(rootGrid, true)

	a.width, a.height = a.ui.GetScreenSize()
	err := a.ui.Run()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initializeHeader() *cview.Flex {
	logo := initializeLogo()
	emptyBar := cview.NewTextView()

	a.statusBar = cview.NewTextView()
	a.statusBar.SetTextAlign(cview.AlignRight)

	flex := cview.NewFlex()
	flex.SetDirection(cview.FlexColumn)
	flex.AddItem(logo, 13, 1, false)
	flex.AddItem(emptyBar, 0, 6, false)
	flex.AddItem(a.statusBar, 13, 1, false)

	return flex
}

func initializeLogo() *cview.TextView {
	textView := cview.NewTextView()
	textView.SetDynamicColors(true)
	textView.SetTextAlign(cview.AlignLeft)

	textView.SetText("[:#ff6600:b]Y[-:-:] Hacker News")
	return textView
}

func (a *App) initializePanels() *cview.Panels {
	if err := a.ui.Init(); err != nil {
		log.Print(err)
	}
	browserPanel := a.initializeTabbedListsLayout()
	postPanel := a.initializeItemContentLayout()

	panels := cview.NewPanels()
	panels.SetBackgroundColor(cview.Styles.PrimitiveBackgroundColor)
	panels.SetBackgroundTransparent(false)

	panels.AddPanel(POSTPANEL, postPanel, true, false)
	panels.AddPanel(LISTPANEL, browserPanel, true, true)

	panels.SetChangedFunc(a.changedPanelsHandler)

	a.initializeInputHandler(browserPanel, postPanel)

	return panels
}

func (a *App) initializeTabbedListsLayout() *cview.Grid {
	a.listView.initializeTabbedLists()

	grid := cview.NewGrid()
	grid.SetBackgroundTransparent(false)
	grid.SetBackgroundColor(cview.Styles.PrimitiveBackgroundColor)

	grid.SetRows(0, 1)
	grid.SetColumns(0)

	grid.AddItem(a.listView.tabbedLists, 0, 0, 1, 2, 1, 1, true)
	grid.AddItem(a.listView.statusBar, 1, 0, 1, 1, 1, 1, false)

	return grid
}

func (a *App) initializeItemContentLayout() *cview.Grid {
	grid := cview.NewGrid()
	grid.SetBackgroundTransparent(false)
	grid.SetBackgroundColor(cview.Styles.PrimitiveBackgroundColor)

	grid.SetRows(0, 0, 1)
	grid.SetColumns(0)

	grid.AddItem(a.postView.content, 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(a.postView.commentsTree, 1, 0, 1, 1, 0, 0, false)
	grid.AddItem(a.postView.debugBar, 2, 0, 1, 1, 0, 0, false)
	return grid
}

func (a *App) changedPanelsHandler() {
	currentTab := a.listView.tabbedLists.GetCurrentTab()

	lastSelectedItemIndex := a.listView.states[currentTab].lastSelectedItemIndex
	if lastSelectedItemIndex != 0 {
		a.listView.states[currentTab].SetCurrentItem(lastSelectedItemIndex)
		a.ui.QueueUpdateDraw(func() {})
	}
}
func (a *App) resizeHandler(width, height int) {
	oldWidth := a.width
	oldHeight := a.height
	if a.listView.tabbedLists.HasFocus() {
		if oldHeight == height && oldWidth != width {
			a.width = width
			a.listView.resizeListItems(width)
		} else if oldHeight != height {
			a.width = width
			a.listView.populateList()
		}
	}
}
