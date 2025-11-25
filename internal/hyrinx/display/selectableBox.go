package display

import (
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/ausro/hyrinx-launcher/config"
	"github.com/ausro/hyrinx-launcher/internal/hyrinx"
)

type selectableBox struct {
	widget.BaseWidget
	app               *hyrinx.Application
	Img               *canvas.Image
	Label             *widget.Label
	menu              *fyne.Menu
	Name              string
	selected, hovered bool

	// For double-click
	lastTap            time.Time
	doubleClickTimeout time.Duration

	onSelected      func()
	onDoubleClicked func()
}

// Creates a new desktop style widget
func NewSelectableCard(app *hyrinx.Application, icon fyne.Resource) *selectableBox {
	card := &selectableBox{
		app:                app,
		Img:                canvas.NewImageFromResource(icon),
		Label:              widget.NewLabelWithStyle(app.Name, fyne.TextAlignCenter, fyne.TextStyle{}),
		Name:               app.Name,
		doubleClickTimeout: 300 * time.Millisecond,
		onSelected:         nil,
		onDoubleClicked:    func() { go hyrinx.Launch(app.Path, app.Options...) },
	}

	// Image settings
	setImageDetails(card.Img)

	// Label settings
	card.Label.Wrapping = fyne.TextWrapWord

	// Menu
	card.menu = card.createMenu()

	// Set selection callback after card is initialized
	card.onSelected = func() { onWidgetSelected(card) }

	card.ExtendBaseWidget(card)
	return card
}

func (s *selectableBox) Tapped(_ *fyne.PointEvent) {
	now := time.Now()

	// Doubleclick handling, as DoubleClicked interface increases response time by ~1 second
	if !s.lastTap.IsZero() && now.Sub(s.lastTap) <= s.doubleClickTimeout {
		if s.onDoubleClicked != nil {
			s.onDoubleClicked()
		}
		return
	}

	s.lastTap = now
	// Toggle selection and notify manager
	s.selected = !s.selected
	s.Refresh()

	if s.onSelected != nil {
		s.onSelected()
	}
}

// setSelected sets the selection state directly (used by selection manager)
func (s *selectableBox) setSelected(selected bool) {
	s.selected = selected
	s.Refresh()
}

func (s *selectableBox) createMenu() *fyne.Menu {
	i := fyne.NewMenuItem("Edit", func() {
		EditAppDialog(*hyrinx.GetRootWindow(), true, &AppDetails{
			Name:  s.Name,
			Path:  s.app.Path,
			Image: "",
			Opts:  strings.Join(s.app.Options, " "),
		}, func(details *AppDetails) {
			s.editDetails(details.Name, details.Path, details.Image)
			s.app.EditOptions(strings.Split(details.Opts, " "))
		})
	})
	d := fyne.NewMenuItem("Delete", func() {
		dialog.ShowConfirm("Delete app?", "Are you sure you'd like to delete this app?", func(ok bool) {
			if !ok {
				return
			}
			removeGridItem(s)
		}, *hyrinx.GetRootWindow())
	})
	m := fyne.NewMenu("Actions", i, d)

	return m
}

func (s *selectableBox) editDetails(name string, path string, image string) {
	if image != "" {
		i, err := fyne.LoadResourceFromPath(image)
		if err != nil {
			// Skip setting img
		} else {
			img := canvas.NewImageFromResource(i)
			setImageDetails(img)
			s.Img = img
		}
	}
	s.Name = name
	s.Label.SetText(s.Name)
	s.Refresh()

	s.app.EditDetails(name, path, image)
	config.CONF.GetCurrentProfile().Update(config.CONF.GetCurrentProfile().Name, GetApplicationsMap())
}

func setImageDetails(img *canvas.Image) {
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(64, 64))
}

func (b *selectableBox) MouseIn(*desktop.MouseEvent) {
	b.hovered = true
	b.Refresh()
}

func (b *selectableBox) MouseMoved(*desktop.MouseEvent) {}

func (b *selectableBox) MouseOut() {
	b.hovered = false
	b.Refresh()
}

func (b *selectableBox) TappedSecondary(e *fyne.PointEvent) {
	widget.ShowPopUpMenuAtPosition(b.menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
}

type selectableBoxRenderer struct {
	box     *selectableBox
	bg      *canvas.Rectangle
	content *fyne.Container
}

func (r *selectableBoxRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
	r.content.Resize(size)
	r.box.Resize(size)
}

func (r *selectableBoxRenderer) MinSize() fyne.Size {
	return fyne.NewSize(r.box.Size().Width, r.content.MinSize().Height)
}

func (r *selectableBoxRenderer) Refresh() {
	// r.content.Objects[0] = r.box.Img

	if r.box.hovered {
		r.bg.FillColor = theme.Color(theme.ColorNameHover)
	} else if r.box.selected {
		r.bg.FillColor = theme.Color(theme.ColorNameSelection)
	} else {
		r.bg.FillColor = theme.Color(theme.ColorNameBackground)
	}

	if r.box.hovered && r.box.selected {
		r.bg.FillColor = theme.Color(theme.ColorNameSelection)
	}

	newSize := fyne.NewSize(r.box.Size().Width, r.content.MinSize().Height)
	r.Layout(newSize)
	canvas.Refresh(r.box)
}

func (r *selectableBoxRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.bg, r.content}
}

func (r *selectableBoxRenderer) Destroy() {}

func (s *selectableBox) CreateRenderer() fyne.WidgetRenderer {
	content := container.NewVBox(
		s.Img,
		s.Label,
	)
	bg := canvas.NewRectangle(color.NRGBA{0, 0, 0, 0})
	r := &selectableBoxRenderer{box: s, bg: bg, content: content}

	return r
}
