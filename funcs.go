package main

import (
	"fmt"
	"log/slog"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	ticker        = time.NewTicker(time.Millisecond)
	atomicCounter uint64
)

func MakeRequest(resChannel chan<- *ChannelInformation) {
	resChannel <- &ChannelInformation{
		Error:     nil,
		IsLoading: true,
	}

	atomicCounter = 0
	ticker.Reset(time.Millisecond)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			atomic.AddUint64(&atomicCounter, 1)
			app.Draw()
		}
	}()

	resp, err := MakeHTTPCall()
	channelInfo := &ChannelInformation{IsLoading: false}

	if err != nil {
		channelInfo.Error = err
		resChannel <- channelInfo
		return
	}

	channelInfo.Response = resp
	resChannel <- channelInfo
}

// ! for rendering
func RenderResponse(resChannel <-chan *ChannelInformation, bodyView *tview.TextView, headerView *tview.Flex) {
	httpMethodText := CreateNewDynamicTextView()

	counterSec := tview.NewBox().SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		if ticker == nil {
			return 0, 0, 0, 0
		}

		floating := time.Duration(atomicCounter) * time.Millisecond
		seconds := strconv.FormatFloat(floating.Abs().Seconds(), 'f', 2, 64)
		tview.Print(screen, fmt.Sprintf("Time: %ss", seconds), x, height/2, width, tview.AlignCenter, tcell.Color100)
		return 0, 0, 0, 0
	})

	headerView.SetDirection(tview.FlexColumn).
		AddItem(httpMethodText, 0, 1, false).
		AddItem(counterSec, 0, 1, false)

	for {

		res, ok := <-resChannel
		slog.Debug("response of the channel: ", "Response", res)

		if !ok {
			slog.Warn("channel was closed")
			break
		}

		httpMethodText.SetText("")
		bodyView.SetText("")

		// header control
		if res.Response != nil {
			ticker.Stop()
			code := res.Response.StatusCode
			fmt.Fprintf(bodyView, "[white]%s[white]", res.Response.ParsedBody)
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

			fmt.Fprintf(httpMethodText, "[%s::b]%s[white]", statusColor, res.Response.Status)
		}

		// body control
		if res.IsLoading {
			fmt.Fprintf(bodyView, "[white]%s[white]", "Loading...")
		} else if res.Error != nil {
			fmt.Fprintf(bodyView, "[red]%s[white]", res.Error)
		}
		app.Draw()
	}
}
