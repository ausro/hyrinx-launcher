package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/ausro/hyrinx-launcher/config"
	"github.com/ausro/hyrinx-launcher/internal/hyrinx/display"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hyrinx Launcher")
	config.InitConfig()

	w.CenterOnScreen()

	al := display.CreateAppLayout()

	w.SetContent(al)
	w.SetMainMenu(display.MakeMainMenu())
	w.SetOnDropped(display.AcceptDropItem())

	// TODO: Set based on config
	w.Resize(fyne.NewSize(1280, 720))
	w.SetFixedSize(false)

	w.ShowAndRun()
}
