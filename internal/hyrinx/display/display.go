package display

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ausro/hyrinx-launcher/config"
	"github.com/ausro/hyrinx-launcher/internal/hyrinx"
)

var appGrid *fyne.Container
var selectedWidget *selectableBox
var widgets []*selectableBox

func CreateAppLayout() *fyne.Container {
	grid := container.New(NewDynamicGridWrapLayout(fyne.NewSize(config.CONF.AppSize, config.CONF.AppSize)))
	gridWrapper := CreateGrid(grid)
	header := CreateHeader()
	overall := container.NewBorder(header, nil, nil, nil, gridWrapper)

	return overall
}

// Creates the core grid layout that contains all applications.
// Returns the clickable wrapper widget to be set as window content.
//   - grid: The fyne container to use
func CreateGrid(grid *fyne.Container) fyne.CanvasObject {
	appGrid = grid
	wrapper := NewClickableGridWrapper(grid)

	go fyne.Do(func() { AddGridItems(config.CONF.GetCurrentProfile().Applications) })
	// AddGridItems(config.CONF.GetCurrentProfile().Applications)

	return wrapper
}

func CreateHeader() fyne.CanvasObject {
	addButton := widget.NewButtonWithIcon("", theme.Icon(theme.IconNameContentAdd), DefaultAddAppDialog())
	h := container.NewHBox(addButton)

	return h
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

func DefaultAddAppDialog() func() {
	return func() {
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
	}
}

func MakeMainMenu() *fyne.MainMenu {
	i := fyne.NewMenuItem("New App", DefaultAddAppDialog())

	// Label is wonky when too small
	s := fyne.NewMenuItem("Small", func() { SetAppSize(0) })
	d := fyne.NewMenuItem("Default", func() { SetAppSize(1) })
	l := fyne.NewMenuItem("Large", func() { SetAppSize(2) })

	rm := fyne.NewMenuItem("Resize Apps", nil)
	rm.ChildMenu = fyne.NewMenu("Size", s, d, l)
	m := fyne.NewMenu("File", i, rm)
	mm := fyne.NewMainMenu(m)

	return mm
}

func SetAppSize(sIndex int) {
	var size float32
	switch sIndex {
	case 0:
		size = 75
	case 1:
		size = 100
	case 2:
		size = 200
	default:
		size = 100
	}
	ResizeGridItems(size)
	config.CONF.AppSize = size
	config.WriteConfig(config.CONF)
}

func AddGridItems(apps []hyrinx.Application) {
	for _, v := range apps {
		addFromConfig(&v)
	}
}

// Adds from config, skipping write/refresh steps
func addFromConfig(app *hyrinx.Application) {
	w := createWidget(app)
	appGrid.Add(w)
	// Update list
	if sb, ok := w.(*selectableBox); ok {
		widgets = append(widgets, sb)
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

	config.CONF.GetCurrentProfile().Update(config.CONF.GetCurrentProfile().Name, GetApplicationsMap())
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

	config.CONF.GetCurrentProfile().Update(config.CONF.GetCurrentProfile().Name, GetApplicationsMap())
}

func ResizeGridItems(size float32) {
	if l, ok := appGrid.Layout.(*dynamicGridWrapLayout); ok {
		l.CellSize = fyne.NewSize(size, size)
	}

	for _, v := range widgets {
		v.Img.SetMinSize(fyne.NewSize(0.64*size, 0.64*size))
	}
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
func GetApplicationsMap() []hyrinx.Application {
	m := []hyrinx.Application{}
	for _, sb := range widgets {
		if sb != nil && sb.app != nil {
			m = append(m, *sb.app)
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
