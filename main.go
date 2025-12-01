package main

import (
	"log/slog"

	"github.com/Aragon-Joaquin/curlWrapper_CLI/logger"
	"github.com/Aragon-Joaquin/curlWrapper_CLI/ui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Application struct {
	tv *tview.Application
	UI *ui.UIService
}

func main() {
	err := logger.Load(logger.DefaultPath(), slog.LevelDebug)

	if err != nil {
		panic(err)
	}

	drawView := tview.NewApplication()
	app := &Application{
		tv: drawView,
		UI: ui.CreateUIService(drawView),
	}

	tview.Styles = tview.Theme{
		TitleColor: tcell.ColorChartreuse,
	}

	mainBox := app.UI.GenerateLayout()
	if err := app.tv.SetRoot(mainBox, true).SetFocus(mainBox).Run(); err != nil {
		slog.Error(err.Error())
	}
}
