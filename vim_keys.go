package main

import (
	"log/slog"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type VIM_KEYS rune

var (
	VIM_KEY_UP    VIM_KEYS = 'j'
	VIM_KEY_DOWN  VIM_KEYS = 'k'
	VIM_KEY_LEFT  VIM_KEYS = 'h'
	VIM_KEY_RIGHT VIM_KEYS = 'l'
)

var (
	currentFocus  = 0
	navigationBox = tview.NewBox()
)

type borderColorInterface interface {
	tview.Primitive
	SetBorderColor(color tcell.Color) *tview.Box
}

func SetVIMNavigationKeys(app *tview.Application, keyMapText *tview.TextView) {
	sliceOfItems := []tview.Primitive{form, respHeader, reqBody, bodyView}
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		slog.Debug("key pressed", "Key", event.Rune(), "KeyString", event.Key())

		if event.Key() == tcell.KeyESC && !navigationBox.HasFocus() {
			app.SetFocus(navigationBox) // loose focus of everything
			ChangeKeyMapMode(FOCUS_KEY_HELP, keyMapText)

			return event
		}

		if !navigationBox.HasFocus() {
			reqBody.SetDisabled(false)
			form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				return event
			})
			return event
		}

		prevBox := sliceOfItems[currentFocus].(borderColorInterface)

		if navigationBox.HasFocus() && event.Key() == tcell.KeyEnter {
			navigationBox.Blur()
			ChangeKeyMapMode(BLUR_KEY_HELP, keyMapText)
			prevBox.SetBorderColor(tcell.ColorDimGray)
			app.SetFocus(sliceOfItems[currentFocus])
			return event
		}

		form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			return nil
		})
		reqBody.SetDisabled(true)

		switch event.Rune() {
		case rune(VIM_KEY_DOWN):
			if currentFocus+2 < len(sliceOfItems) {
				currentFocus += 2
			}

		case rune(VIM_KEY_UP):
			if currentFocus-2 >= 0 {
				currentFocus -= 2
			}

		case rune(VIM_KEY_LEFT):
			if currentFocus-1 >= 0 {
				currentFocus -= 1
			}

		case rune(VIM_KEY_RIGHT):
			if currentFocus+1 < len(sliceOfItems) {
				currentFocus += 1
			}

		case 'q':
			app.Stop()

		default:
			return event
		}

		prevBox.SetBorderColor(tcell.ColorDarkSlateGray)

		if currentFocus >= 0 && currentFocus < len(sliceOfItems) {
			box := sliceOfItems[currentFocus].(borderColorInterface)
			box.SetBorderColor(tcell.ColorRed)
			app.SetFocus(box)
		}

		return event
	})
}
