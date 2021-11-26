package ui

import (
	"fmt"
	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
	"strings"
)

const (
	logo = `
 **      **                   **                    ****     **                           
/**     /**                  /**                   /**/**   /**                           
/**     /**  ******    ***** /**  **  *****  ******/**//**  /**  *****  ***     **  ******
/********** //////**  **///**/** **  **///**//**//*/** //** /** **///**//**  * /** **//// 
/**//////**  ******* /**  // /****  /******* /** / /**  //**/**/******* /** ***/**//***** 
/**     /** **////** /**   **/**/** /**////  /**   /**   //****/**////  /****/**** /////**
/**     /**//********//***** /**//**//******/***   /**    //***//****** ***/ ///** ****** 
//      //  ////////  /////  //  //  ////// ///    //      ///  ////// ///    /// //////  
`
	subtitle = "Terminal-based HackerNews Client"
)

func Cover() *cview.Flex {
	lines := strings.Split(logo, "\n")
	logoWidth := 0
	logoHeight := len(lines)

	for _, line := range lines {
		if len(line) > logoWidth {
			logoWidth = len(line)
		}
	}

	logoBox := cview.NewTextView()
	logoBox.SetTextColor(Orange.TrueColor())

	fmt.Fprint(logoBox, logo)

	frame := cview.NewFrame(cview.NewBox())
	frame.SetBorders(0, 0, 0, 0, 0, 0)
	frame.AddText(subtitle, true, cview.AlignCenter, tcell.ColorDarkMagenta.TrueColor())

	subFlex := cview.NewFlex()
	subFlex.AddItem(cview.NewBox(), 0, 1, false)
	subFlex.AddItem(logoBox, logoWidth, 1, true)
	subFlex.AddItem(cview.NewBox(), 0, 1, false)

	flex := cview.NewFlex()
	flex.SetDirection(cview.FlexRow)
	flex.AddItem(cview.NewBox(), 0, 7, false)
	flex.AddItem(subFlex, logoHeight, 1, true)
	flex.AddItem(frame, 0, 10, false)
	return flex
}
