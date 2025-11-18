package display

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type gridWrapper struct {
	widget.BaseWidget
	grid *fyne.Container
}

func NewClickableGridWrapper(grid *fyne.Container) *gridWrapper {
	w := &gridWrapper{
		grid: grid,
	}
	w.ExtendBaseWidget(w)
	return w
}

func (w *gridWrapper) Tapped(_ *fyne.PointEvent) {
	// Deselect item when empty space is tapped
	deselectAll()
}

func (w *gridWrapper) CreateRenderer() fyne.WidgetRenderer {
	return &gridWrapperRenderer{
		wrapper: w,
	}
}

type gridWrapperRenderer struct {
	wrapper *gridWrapper
}

func (r *gridWrapperRenderer) Layout(size fyne.Size) {
	r.wrapper.grid.Resize(size)
}

func (r *gridWrapperRenderer) MinSize() fyne.Size {
	return r.wrapper.grid.MinSize()
}

func (r *gridWrapperRenderer) Refresh() {
	canvas.Refresh(r.wrapper.grid)
}

func (r *gridWrapperRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.wrapper.grid}
}

func (r *gridWrapperRenderer) Destroy() {}
