package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Aragon-Joaquin/curlWrapper_CLI/ui"
	ut "github.com/Aragon-Joaquin/curlWrapper_CLI/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FieldState struct {
	URLField    string
	MethodField ut.HTTPMethod
	Body        string
}

type ChannelInformation struct {
	Response  *RequestJson
	Error     error
	IsLoading bool
}

var GlobalFieldState = &FieldState{}
var app *tview.Application

func main() {
	resChannel := make(chan *ChannelInformation)

	//! creating the form
	urlInput := ui.CreateURLInput().
		SetChangedFunc(func(text string) {
			slog.Debug("Input Text Changed: ", "TEXT", text)
			GlobalFieldState.URLField = text
		})

	methodDMenu := ui.CreateDropDownMethods().
		SetSelectedFunc(func(text string, index int) {
			slog.Debug("DropDown Changed: ", "TEXT", text, "INDEX", index)
			GlobalFieldState.MethodField = ut.HTTPMethod(index)
		})

	form := ui.CreateFormMenu(urlInput, methodDMenu).
		AddButton("Make Request", func() {
			app.SetTitle("Done!")

			time.AfterFunc(time.Second*1, func() {
				app.SetTitle(ut.APP_NAME)
			})

			go MakeRequest(resChannel)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle(ut.APP_NAME)

	//reqBody
	jsonIsValid := ui.CreateNewDynamicTextView()
	keyMap := ui.CreateNewDynamicTextView()
	reqBody := ui.CreateBodyInput()

	reqBody.SetBorder(true).SetTitle("Write the JSON here")

	jsonIsValid.SetText("Waiting for input...")
	jsonIsValid.SetTextAlign(tview.AlignCenter)

	keyMap.SetText("<ESC> to Unfocus")
	keyMap.SetTextAlign(tview.AlignCenter)

	reqBody.SetChangedFunc(func() {
		//todo: add debouncer here
		text := reqBody.GetText()

		jsonIsValid.SetText("")
		if len(text) < 2 {
			fmt.Fprint(jsonIsValid, "[white]Waiting for input...[white]")
			return
		}

		slog.Debug("TextArea is JSON: ", "Bool", IsJSON(text))
		if IsJSON(text) {
			GlobalFieldState.Body = text
			fmt.Fprint(jsonIsValid, "[green]Valid JSON[white]")
		} else {
			GlobalFieldState.Body = ""
			fmt.Fprint(jsonIsValid, "[red]Invalid JSON[white]")
		}
	})

	reqBody.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC {
			reqBody.Blur()
		}
		return event
	})

	reqBodyUtils := tview.NewGrid().
		SetRows(0, 1).
		AddItem(reqBody, 0, 0, 1, 2, 0, 0, true).
		AddItem(jsonIsValid, 1, 0, 1, 1, 0, 0, false).
		AddItem(keyMap, 1, 1, 1, 1, 0, 0, false)

	flexRequest := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true).
		AddItem(reqBodyUtils, 0, 5, false)

	//! json response layout
	headerView := ui.CreateNewDynamicTextView()
	bodyView := ui.CreateNewDynamicTextView()

	headerView.SetBorder(true).SetTitle("Response Header")
	bodyView.SetBorder(true).SetTitle("Response Body")

	go RenderResponse(resChannel, bodyView, headerView)

	flexResponse := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(headerView, 0, 1, true).
		AddItem(bodyView, 0, 10, true)

	// final layout
	mainBox := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(flexRequest, 0, 1, true).
		AddItem(flexResponse, 0, 1, false)

	if err := app.SetRoot(mainBox, true).Run(); err != nil {
		slog.Error(err.Error())
	}
}

// ? init funcs
func init() {
	err := LoggerLoad(LoggerDefaultPath(), slog.LevelDebug)

	if err != nil {
		panic(err)
	}
}

func init() {
	newApp := tview.NewApplication()
	newApp.EnableMouse(true)
	newApp.SetTitle(ut.APP_NAME)

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
