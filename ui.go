package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateURLInput() *tview.InputField {
	return tview.NewInputField().
		SetLabel("URL: ").
		SetPlaceholder(DEFAULT_URL_EXAMPLE).
		SetFieldWidth(30).
		SetAcceptanceFunc(nil)
}

func CreateDropDownMethods() *tview.DropDown {
	return tview.NewDropDown().
		SetLabel("HTTP Method: ").
		SetOptions(ALL_HTTP_METHODS, func(text string, index int) {}).
		SetCurrentOption(0)
}

func AddItemsFormMenu(form *tview.Form, items ...tview.FormItem) *tview.Form {
	for _, box := range items {
		form.AddFormItem(box)
	}

	return form
}

func CreateNewDynamicTextView() *tview.TextView {
	return tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
}

func CreateBodyInput() *tview.TextArea {
	return tview.NewTextArea().
		SetPlaceholder(EXAMPLE_JSON)
}

type KEYMAP_MODES uint8

const (
	BLUR_KEY_HELP KEYMAP_MODES = iota
	FOCUS_KEY_HELP
)

var KEYMAP_TEXT = map[KEYMAP_MODES]string{
	BLUR_KEY_HELP:  "<ESC> to Blur",
	FOCUS_KEY_HELP: "<ENTER> to Focus",
}

func ChangeKeyMapMode(mode KEYMAP_MODES, keyMap *tview.TextView) {
	text := KEYMAP_TEXT[mode]

	keyMap.SetText(text)

	switch mode {
	case FOCUS_KEY_HELP:
		keyMap.SetTextColor(tcell.ColorDarkOrange)

	case BLUR_KEY_HELP:
		keyMap.SetTextColor(tcell.ColorDarkCyan)
	}
}
