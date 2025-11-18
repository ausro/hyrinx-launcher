package display

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/ausro/hyrinx-launcher/config"
	"github.com/ausro/hyrinx-launcher/internal/hyrinx"
)

var appGrid *fyne.Container
var selectedWidget *selectableBox
var widgets []*selectableBox

// Creates the core grid layout that contains all applications.
// Returns the clickable wrapper widget to be set as window content.
//   - grid: The fyne container to use
func Create(grid *fyne.Container) fyne.CanvasObject {
	appGrid = grid
	wrapper := NewClickableGridWrapper(grid)

	//TESTING REMOVE BEFORE DEPLOYMENT
	//addGridItem(hyrinx.CreateApplication("banana", "C:/Users/Overseer/AppData/Local/Programs/Anki/anki.exe", "icon.png"))
	go AddGridItems(config.CONF.GetCurrentProfile().Applications)

	return wrapper
}

// Deselect all widgets and clear the selection tracker
func deselectAll() {
	if selectedWidget != nil {
		selectedWidget.setSelected(false)
		selectedWidget = nil
	}
}

// Notify the manager that a widget was selected
func onWidgetSelected(w *selectableBox) {
	if selectedWidget != w {
		if selectedWidget != nil {
			selectedWidget.setSelected(false)
		}
		selectedWidget = w
	}
}

func MakeMainMenu() *fyne.MainMenu {
	i := fyne.NewMenuItem("New App", func() {
		EditAppDialog(*hyrinx.GetRootWindow(), true, &AppDetails{}, func(details *AppDetails) {
			options := strings.Split(details.Opts, " ")
			a := &hyrinx.Application{
				Name:    details.Name,
				Path:    details.Path,
				Icon:    details.Image,
				Options: options,
			}

			addGridItem(a)
		})
	})
	m := fyne.NewMenu("File", i)
	mm := fyne.NewMainMenu(m)

	return mm
}

func AddGridItems(apps map[int]hyrinx.Application) {
	for _, v := range apps {
		addGridItem(&v)
	}
}

// Add a new application to the grid
//   - app: The application to add
func addGridItem(app *hyrinx.Application) {
	w := createWidget(app)
	appGrid.Add(w)
	// Update list
	if sb, ok := w.(*selectableBox); ok {
		widgets = append(widgets, sb)
	}
	appGrid.Refresh()

	config.CONF.GetCurrentProfile().Update("", GetApplicationsMap())
}

// Removes an application from the grid
//   - w: The widget that was targeted
func removeGridItem(w fyne.CanvasObject) {
	appGrid.Remove(w)
	appGrid.Refresh()

	// Remove from list
	if sb, ok := w.(*selectableBox); ok {
		for i, v := range widgets {
			if v == sb {
				widgets = append(widgets[:i], widgets[i+1:]...)
				break
			}
		}
	}

	config.CONF.GetCurrentProfile().Update("", GetApplicationsMap())
}

func createWidget(app *hyrinx.Application) fyne.CanvasObject {
	var icon fyne.Resource
	var err error
	if app.Icon != "" {
		icon, err = fyne.LoadResourceFromPath(app.Icon)
		if err != nil {
			icon = theme.BrokenImageIcon()
		}
	} else {
		icon = theme.FileApplicationIcon()
	}

	w := NewSelectableCard(app, icon)

	return w
}

// Get the current applications slice as a map
func GetApplicationsMap() map[int]hyrinx.Application {
	m := make(map[int]hyrinx.Application)
	for i, sb := range widgets {
		if sb != nil && sb.app != nil {
			m[i] = *sb.app
		}
	}
	return m
}

// Drag-drop handler to create an application from the dropped items
func AcceptDropItem() func(fyne.Position, []fyne.URI) {
	return func(p fyne.Position, u []fyne.URI) {
		for _, v := range u {
			a := hyrinx.CreateApplication(v.Name(), v.Path(), "")
			addGridItem(a)
		}
	}
}
