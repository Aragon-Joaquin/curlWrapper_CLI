package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// modules
var (
	form       *tview.Form
	respHeader *tview.Flex
	reqBody    *tview.TextArea
	bodyView   *tview.TextView
)

type ChannelInformation struct {
	Response  *RequestJson
	Error     error
	IsLoading bool
}

// ? init funcs
func init() {
	err := LoggerLoad(LoggerDefaultPath(), slog.LevelDebug)
	if err != nil {
		panic(err)
	}
}

func init() {
	respHeader = tview.NewFlex()
	bodyView = CreateNewDynamicTextView()
	reqBody = CreateBodyInput()
	form = tview.NewForm().
		SetButtonsAlign(tview.AlignCenter)
}

// ? entry point
func main() {
	resChannel := make(chan *ChannelInformation)

	//! creating the form
	urlInput := CreateURLInput().
		SetChangedFunc(func(text string) {
			slog.Debug("Input Text Changed: ", "TEXT", text)
			GlobalFieldState.URLField = text
		})

	methodDMenu := CreateDropDownMethods().
		SetSelectedFunc(func(text string, index int) {
			slog.Debug("DropDown Changed: ", "TEXT", text, "INDEX", index)
			GlobalFieldState.MethodField = HTTPMethod(index)
		})

	AddItemsFormMenu(form, urlInput, methodDMenu).
		AddButton("Make Request", func() {
			app.SetTitle("Done!")

			time.AfterFunc(time.Second*1, func() {
				app.SetTitle(APP_NAME)
			})

			go MakeRequest(resChannel)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle(APP_NAME).SetBorderColor(tcell.ColorDimGray)

	// reqBody
	jsonIsValid := CreateNewDynamicTextView()
	keyMap := CreateNewDynamicTextView()

	reqBody.SetBorder(true).SetTitle("Write the JSON here")

	jsonIsValid.SetText("Waiting for input...")
	jsonIsValid.SetTextAlign(tview.AlignCenter)

	keyMap.SetTextAlign(tview.AlignCenter)
	ChangeKeyMapMode(BLUR_KEY_HELP, keyMap)

	reqBody.SetChangedFunc(func() {
		// todo: add debouncer here
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

	reqBodyUtils := tview.NewGrid().
		SetRows(0, 1).
		AddItem(reqBody, 0, 0, 1, 2, 0, 0, true).
		AddItem(jsonIsValid, 1, 0, 1, 1, 0, 0, false).
		AddItem(keyMap, 1, 1, 1, 1, 0, 0, false)

	flexRequest := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 2, true).
		AddItem(reqBodyUtils, 0, 5, false)

		//! json response layout

	respHeader.SetDirection(tview.FlexColumn).SetBorder(true).SetTitle("Response Header")
	bodyView.SetBorder(true).SetTitle("Response Body")

	go RenderResponse(resChannel, bodyView, respHeader)

	flexResponse := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(respHeader, 0, 1, true).
		AddItem(bodyView, 0, 10, true)

	// final layout
	mainBox := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(flexRequest, 0, 1, true).
		AddItem(flexResponse, 0, 1, false)

	//? extras
	SetVIMNavigationKeys(app, keyMap)

	if err := app.SetRoot(mainBox, true).Run(); err != nil {
		slog.Error(err.Error())
	}
}
