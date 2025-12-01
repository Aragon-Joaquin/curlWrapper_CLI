package ui

import "github.com/rivo/tview"

func (t *UIService) CreateInput(name, placeholder string) *tview.InputField {
	input := tview.NewInputField()

	input.SetPlaceholder(placeholder)
	input.SetLabel(name)

	return input
}
