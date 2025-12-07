package ui

import "github.com/rivo/tview"

type UIService struct {
	*tview.Application
}

func CreateUIService(t *tview.Application) *UIService {
	return &UIService{t}
}
