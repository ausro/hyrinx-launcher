package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/ausro/hyrinx-launcher/config"
	"github.com/ausro/hyrinx-launcher/internal/hyrinx/display"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hyrinx Launcher")
	config.InitConfig()

	w.CenterOnScreen()
	grid := container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 100)))
	gridWrapper := display.Create(grid)

	w.SetContent(gridWrapper)
	w.SetMainMenu(display.MakeMainMenu())
	w.SetOnDropped(display.AcceptDropItem())

	// TODO: Set based on config
	w.Resize(fyne.NewSize(1280, 720))
	w.SetFixedSize(false)

	w.ShowAndRun()
}
