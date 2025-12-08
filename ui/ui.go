package ui

import (
	"github.com/Aragon-Joaquin/curlWrapper_CLI/utils"
	"github.com/rivo/tview"
)

func CreateURLInput() *tview.InputField {
	return tview.NewInputField().
		SetLabel("URL: ").
		SetPlaceholder(utils.DEFAULT_URL_EXAMPLE).
		SetFieldWidth(30).
		SetAcceptanceFunc(nil)
}

func CreateDropDownMethods() *tview.DropDown {
	return tview.NewDropDown().
		SetLabel("HTTP Method: ").
		SetOptions(utils.ALL_HTTP_METHODS, func(text string, index int) {}).
		SetCurrentOption(0)
}

func CreateFormMenu(items ...tview.FormItem) *tview.Form {
	form := tview.NewForm().
		SetButtonsAlign(tview.AlignCenter)

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
		SetPlaceholder(utils.EXAMPLE_JSON)
}
