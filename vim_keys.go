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
	currentFocus  = 0 // -1 == nothing focused
	navigationBox = tview.NewBox()
)

func SetVIMNavigationKeys(app *tview.Application) {
	sliceOfItems := []tview.Primitive{form, headerView, reqBody, bodyView}
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		slog.Debug("key pressed", "Key", event.Rune())
		switch event.Key() {
		case tcell.KeyESC:
			app.SetFocus(navigationBox) //loose focus of everything
			currentFocus = -1
			return event
		}

		if !navigationBox.HasFocus() {
			reqBody.SetDisabled(false)
			form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				return event
			})
			return event
		}

		if navigationBox.HasFocus() && event.Key() == tcell.KeyEnter {
			navigationBox.Blur()
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
				app.SetFocus(sliceOfItems[currentFocus])
			}

		case rune(VIM_KEY_UP):
			if currentFocus-2 >= 0 {
				currentFocus -= 2
				app.SetFocus(sliceOfItems[currentFocus])
			}

		case rune(VIM_KEY_LEFT):
			if currentFocus-1 >= 0 {
				currentFocus -= 1
				app.SetFocus(sliceOfItems[currentFocus])
			}

		case rune(VIM_KEY_RIGHT):
			if currentFocus+1 < len(sliceOfItems) {
				currentFocus += 1
				app.SetFocus(sliceOfItems[currentFocus])
			}
		}

		return event

	})
}
