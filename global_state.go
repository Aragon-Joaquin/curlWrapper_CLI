package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FieldState struct {
	URLField    string
	MethodField HTTPMethod
	Body        string
}

var GlobalFieldState = &FieldState{}
var app *tview.Application

func init() {
	newApp := tview.NewApplication()
	newApp.EnableMouse(true)
	newApp.SetTitle(APP_NAME)

	tview.Styles = tview.Theme{
		// input
		ContrastBackgroundColor:    tcell.ColorDarkRed,    //bg-color
		PrimaryTextColor:           tcell.ColorWhiteSmoke, //text-color
		ContrastSecondaryTextColor: tcell.ColorDimGray,    // placeholder-color
		SecondaryTextColor:         tcell.ColorWhiteSmoke, //label-color

		TitleColor:                  tcell.ColorGreenYellow,
		MoreContrastBackgroundColor: tcell.ColorDarkRed, // menus-bg-color
		BorderColor:                 tcell.ColorDarkSlateGray,
	}

	app = newApp
}
