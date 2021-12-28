package ui

import "github.com/gdamore/tcell/v2"

type rgb struct {
	r int
	g int
	b int
}

var (
	Orange       = tcell.NewRGBColor(255, 102, 0)
	ScoreHeatMap = []string{
		"#cc3d00",
		"#d02b27",
		"#d01b3f",
		"#ca1555",
		"#bf1c69",
		"#b0297b",
		"#9d3689",
		"#864093",
		"#6e4898",
		"#534e99",
		"#355297",
		"#0c5490",
	}
)
