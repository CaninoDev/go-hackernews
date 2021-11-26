package ui

import 	"code.rocketnine.space/tslocum/cview"

func StatusBar() cview.Primitive {
	statusBar := cview.NewTextView()

	statusBar.SetTextColor(Orange)
	return statusBar
}
