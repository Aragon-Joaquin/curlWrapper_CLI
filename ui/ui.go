package ui

import "github.com/rivo/tview"

type UIService struct {
	*tview.Application
}

func CreateUIService(t *tview.Application) *UIService {
	return &UIService{t}
}

func (t *UIService) GenerateLayout() *tview.Grid {
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	main := newPrimitive("Make the request here")
	side := newPrimitive("JSON Response")

	grid := tview.NewGrid().
		SetRows().
		SetColumns(0, 70).
		SetBorders(true).
		AddItem(main, 0, 0, 1, 1, 0, 0, false).
		AddItem(side, 0, 1, 1, 1, 0, 0, false)

	return grid
}
