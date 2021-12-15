package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/CaninoDev/go-hackernews/internal/store"
	"log"
)

type App struct {
	ui           *cview.Application
	focusManager *cview.FocusManager
	store        *store.Store
	statusBar    *cview.TextView
	panels       *cview.Panels
	Cover        *cview.Flex
	browser      *Browser
	post         *Post
}

func New() (*App, error) {
	app := &App{}
	app.ui = cview.NewApplication()

	store, err := store.New()
	if err != nil {
		return &App{}, err
	}

	app.store = store

	app.browser = NewBrowser(app)
	app.post = NewPost(app)
	app.statusBar = cview.NewTextView()

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

	a.initializeData()

	err := a.ui.Run()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initializeHeader() *cview.Grid {
	logo := initializeLogo()
	emptyBar := cview.NewTextView()
	a.statusBar = cview.NewTextView()

	grid := cview.NewGrid()
	grid.SetBackgroundColor(cview.Styles.PrimitiveBackgroundColor)
	grid.SetBackgroundTransparent(false)
	grid.SetRows(0)
	grid.SetColumns(13, 0, 9)
	grid.AddItem(logo, 1, 1, 1, 1, 1, 0, false)
	grid.AddItem(emptyBar, 1, 2, 1, 1, 1, 0, false)
	grid.AddItem(a.statusBar, 1, 3, 1, 1, 1, 0, false)

	return grid
}

func initializeLogo() *cview.TextView {
	textView := cview.NewTextView()
	textView.SetDynamicColors(true)
	textView.SetText("[:#ff6600:b]Y[-:-:] Hacker News")
	return textView
}

func (a *App) initializePanels() *cview.Panels {
	if err := a.ui.Init(); err != nil {
		log.Print(err)
	}
	browserPanel := a.initializeBrowserLayout()
	postPanel := a.initializePostLayout()

	panels := cview.NewPanels()
	panels.SetBackgroundColor(cview.Styles.PrimitiveBackgroundColor)
	panels.SetBackgroundTransparent(false)
	panels.AddPanel("post", postPanel, true, false)
	panels.AddPanel("browser", browserPanel, true, true)

	a.initializeInputHandler(browserPanel, postPanel)
	return panels
}

func (a *App) initializeBrowserLayout() *cview.Grid {
	a.browser.initializeTabbedLists()

	grid := cview.NewGrid()
	grid.SetBackgroundTransparent(false)
	grid.SetBackgroundColor(cview.Styles.PrimitiveBackgroundColor)
	grid.SetRows(0, 1)
	grid.SetColumns(0, 0)
	grid.AddItem(a.browser.lists, 0, 0, 1, 2, 1, 1, true)
	grid.AddItem(a.browser.debugBar, 1, 0, 1, 1, 1, 1, false)
	grid.AddItem(a.browser.statusBar, 1, 1, 1, 1, 1, 1, false)

	return grid
}

func (a *App) initializePostLayout() *cview.Grid {
	grid := cview.NewGrid()
	grid.SetBackgroundTransparent(false)
	grid.SetBackgroundColor(cview.Styles.PrimitiveBackgroundColor)
	grid.SetRows(0, 1)
	grid.SetColumns(0)
	grid.AddItem(a.post.commentsTree, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(a.post.debugBar, 1, 0, 1, 1, 0, 0, false)

	return grid
}

func (a *App) initializeData() {
	//a.browser.populate()
}
