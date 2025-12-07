package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	APP_NAME            = "Curl Wrapper"
	DEFAULT_URL_EXAMPLE = "http://localhost:5172"
)

type FieldState struct {
	URLField    string
	MethodField HTTPMethod
}

type ChannelInformation struct {
	Response  *RequestJson
	Error     error
	IsLoading bool
}

var GlobalFieldState = &FieldState{}

func main() {
	err := LoggerLoad(LoggerDefaultPath(), slog.LevelDebug)

	if err != nil {
		panic(err)
	}

	resChannel := make(chan *ChannelInformation)
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

	//! creating the form
	urlInput := tview.NewInputField().
		SetLabel("URL: ").
		SetPlaceholder(DEFAULT_URL_EXAMPLE).
		SetFieldWidth(30).
		SetAcceptanceFunc(nil).
		SetChangedFunc(func(text string) {
			slog.Debug("Input Text Changed: ", "TEXT", text)
			GlobalFieldState.URLField = text
		})

	methodDMenu := tview.NewDropDown().
		SetLabel("HTTP Method: ").
		SetOptions(ALL_HTTP_METHODS, func(text string, index int) {
			slog.Debug("DropDown Changed: ", "TEXT", text, "INDEX", index)
			GlobalFieldState.MethodField = HTTPMethod(index)
		}).
		SetCurrentOption(0)

	form := tview.NewForm().
		AddFormItem(urlInput).
		AddFormItem(methodDMenu).
		AddButton("Make Request", func() {
			app.SetTitle("Done!")

			resChannel <- &ChannelInformation{
				Response:  nil,
				Error:     nil,
				IsLoading: true,
			}

			time.AfterFunc(time.Second*1, func() {
				app.SetTitle(APP_NAME)
			})

			go func() {
				resp, err := MakeHTTPCall()

				channelInfo := &ChannelInformation{IsLoading: false}

				if err != nil {
					channelInfo.Error = err
					resChannel <- channelInfo
					return
				}

				channelInfo.Response = resp
				resChannel <- channelInfo

			}()
		}).
		AddButton("Quit", func() {
			app.Stop()
		}).
		SetButtonsAlign(tview.AlignCenter)

	form.SetBorder(true).SetTitle(APP_NAME)

	//! json response layout
	// header
	headerView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)

	// body
	bodyView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)

	go func() {
		for {
			res, ok := <-resChannel

			slog.Debug("response of the channel: ", "Response", res)

			if !ok {
				slog.Warn("channel was closed")
				break
			}

			bodyView.SetText("")
			headerView.SetText("")

			//header control
			if res.Response != nil {
				code := res.Response.StatusCode
				var statusColor string
				switch {
				case code < 200:
					statusColor = "skyblue"
				case code >= 200 && code < 300:
					statusColor = "green"
				case code >= 300 && code < 400:
					statusColor = "violet"
				case code >= 400 && code < 500:
					statusColor = "yellow"
				case code >= 500:
					statusColor = "red"
				}

				fmt.Fprintf(headerView, "[%s::b]%s[white]", statusColor, res.Response.Status)
			}

			//body control
			if res.IsLoading {
				fmt.Fprintf(bodyView, "[white]%s[white]", "Loading...")
			} else if res.Error != nil {
				fmt.Fprintf(bodyView, "[red]%s[white]", res.Error)
			}

			if res.Response != nil {
				fmt.Fprintf(bodyView, "[white]%s[white]", res.Response.ParsedBody)
			}

			app.Draw()
		}
	}()

	headerView.SetBorder(true).SetTitle("Response Header")
	bodyView.SetBorder(true).SetTitle("Response Body")

	flexResponse := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(headerView, 0, 1, true).
		AddItem(bodyView, 0, 10, true)

	// final layout
	mainBox := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(form, 0, 1, true).
		AddItem(flexResponse, 0, 1, false)

	if err := app.SetRoot(mainBox, true).Run(); err != nil {
		slog.Error(err.Error())
	}
}
