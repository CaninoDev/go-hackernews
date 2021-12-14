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
	root         *cview.Panels
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

	return app, nil
}

func (a *App) Start() error {

	a.root = a.initializePanels()
	a.initializeData()
	a.ui.SetInputCapture(a.inputHandler)
	a.ui.SetRoot(a.root, true)

	err := a.ui.Run()
	if err != nil {
		return err
	}
	return nil
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

//func NewStore(str *store.Store) *App {
//	var app App
//
//	app.ui = cview.NewApplication()
//	defer app.ui.HandlePanic()
//
//	app.store = str
//
//	app.root = cview.NewGrid()
//	app.Cover = Cover()
//	app.StatusBar = cview.NewProgressBar()
//	app.DebugBar = cview.NewTextView()
//
//	app.ListPanels = NewListPanels()
//
//	app.CommentsTree = NewComments()
//
//	commentDisplay := cview.NewFlex()
//	commentDisplay.AddItem(app.CommentsTree.Title, 1, 0, false)
//	commentDisplay.AddItem(app.CommentsTree.Tree, 0, 1, false)
//	commentDisplay.AddItem(app.CommentsTree.Text, 0, 1, false)
//
//	mainPane := cview.NewFlex()
//	mainPane.SetDirection(cview.FlexRow)
//	mainPane.AddItem(app.ListPanels, 0, 1, true)
//	mainPane.AddItem(commentDisplay, 0, 1, false)
//
//	app.root = cview.NewGrid()
//	app.root.SetRows(0, 1, 1)
//	app.root.SetColumns(0)
//	app.root.SetBackgroundColor(cview.Styles.PrimitiveBackgroundColor)
//
//	//app.root.AddItem(app.Cover,0,0,1,1,0,0,false)
//	app.root.AddItem(mainPane, 0, 0, 1, 1, 0, 0, false)
//	app.root.AddItem(app.StatusBar, 1, 0, 1, 1, 0, 0, false)
//	app.root.AddItem(app.DebugBar, 2, 0, 1,1,0,0,false )
//	//flex := cview.NewFlex()
//	//flex.SetDirection(cview.FlexRow)
//	//flex.AddItemAtIndex(0, Cover(), 0, 1, false)
//	//flex.AddItemAtIndex(1, app.StatusBar, 1, 0, false)
//
//	app.ui.SetInputCapture(app.inputHandler)
//	app.ui.SetRoot(app.root, true)
//
//	return &app
//}
//
