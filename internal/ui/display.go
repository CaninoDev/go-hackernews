package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"github.com/CaninoDev/go-hackernews/internal/api"
)

type Display struct {
	App       *cview.Application
	DB        *api.FirebaseClient
	Root      *cview.Grid
	Cover    *cview.Flex
	Posts    *Posts
	Comments *Comments
	StatusBar *cview.ProgressBar
	DebugBar *cview.TextView
}

func Init(db *api.FirebaseClient) *Display {
	var app = cview.NewApplication()
	var display Display

	display.App = app
	defer display.App.HandlePanic()

	display.DB = db

	display.Root = cview.NewGrid()
	display.Cover = Cover()
	display.StatusBar = cview.NewProgressBar()
	display.DebugBar = cview.NewTextView()

	display.Posts = display.NewPostsView()

	display.Comments = NewComments()

	commentDisplay := cview.NewFlex()
	commentDisplay.AddItem(display.Comments.Title, 1, 0, false)
	commentDisplay.AddItem(display.Comments.Tree, 0, 1, false)
	commentDisplay.AddItem(display.Comments.Text, 0, 1, false)

	mainPane := cview.NewFlex()
	mainPane.SetDirection(cview.FlexRow)
	mainPane.AddItem(display.Posts, 0, 1, true)
	mainPane.AddItem(commentDisplay, 0, 1, false)

	display.Root = cview.NewGrid()
	display.Root.SetRows(0, 1, 1)
	display.Root.SetColumns(0)
	display.Root.SetBackgroundColor(cview.Styles.PrimitiveBackgroundColor)

	//display.Root.AddItem(display.Cover,0,0,1,1,0,0,false)
	display.Root.AddItem(mainPane, 0, 0, 1, 1, 0, 0, false)
	display.Root.AddItem(display.StatusBar, 1, 0, 1, 1, 0, 0, false)
	display.Root.AddItem(display.DebugBar, 2, 0, 1,1,0,0,false )
	//flex := cview.NewFlex()
	//flex.SetDirection(cview.FlexRow)
	//flex.AddItemAtIndex(0, Cover(), 0, 1, false)
	//flex.AddItemAtIndex(1, display.StatusBar, 1, 0, false)

	display.App.SetInputCapture(display.AppKeyHandler)
	display.App.SetRoot(display.Root, true)

	return &display
}

func (d *Display) Run() error {
	err := d.App.Run()
	if err != nil {
		return err
	}
	return nil
}
