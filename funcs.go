package main

import (
	"fmt"
	"log/slog"

	"github.com/rivo/tview"
)

func MakeRequest(resChannel chan<- *ChannelInformation) {
	resChannel <- &ChannelInformation{
		Response:  nil,
		Error:     nil,
		IsLoading: true,
	}

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
func RenderResponse(resChannel <-chan *ChannelInformation, bodyView, headerView *tview.TextView) {
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

}
