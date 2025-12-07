package main

import (
	"log/slog"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	APP_NAME            = "Curl Wrapper"
	DEFAULT_URL_EXAMPLE = "http://localhost:5172"
)

func main() {
	err := LoggerLoad(LoggerDefaultPath(), slog.LevelDebug)

	if err != nil {
		panic(err)
	}

	slog.Debug("App is running normally")

	//! creating the layout
	app := tview.NewApplication()
	app.EnableMouse(true) // not that im going to use this
	app.SetTitle(APP_NAME)

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

	//creating the form

	urlInput := tview.NewInputField().
		SetLabel("URL: ").
		SetPlaceholder(DEFAULT_URL_EXAMPLE).
		SetFieldWidth(30).
		SetAcceptanceFunc(nil).
		SetChangedFunc(func(text string) {
			slog.Debug("Input Text Changed: ", "TEXT", text)
		})

	methodDMenu := tview.NewDropDown().
		SetLabel("HTTP Method: ").
		SetOptions(ALL_HTTP_METHODS, func(text string, index int) {
			slog.Debug("DropDown Changed: ", "TEXT", text, "INDEX", index)
		}).
		SetCurrentOption(0)

	form := tview.NewForm().
		AddFormItem(urlInput).
		AddFormItem(methodDMenu).
		AddButton("Save", func() {
			app.SetTitle("Saved!")
			time.AfterFunc(time.Second*1, func() {
				app.SetTitle(APP_NAME)
			})
		}).
		AddButton("Quit", func() {
			app.Stop()
		}).
		SetButtonsAlign(tview.AlignCenter)

	form.SetBorder(true).SetTitle(APP_NAME)
	// final layout
	mainBox := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(form, 0, 1, true).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("JSON Response"), 0, 1, false)

	if err := app.SetRoot(mainBox, true).Run(); err != nil {
		slog.Error(err.Error())
	}
}
