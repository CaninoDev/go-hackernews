package ui

import 	"code.rocketnine.space/tslocum/cview"

type StatusBar struct {
	*cview.ProgressBar
}

func NewStatusBar() *StatusBar {
	return &StatusBar{
		ProgressBar: cview.NewProgressBar(),
	}
}
