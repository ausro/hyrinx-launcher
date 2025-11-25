package display

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type dynamicGridWrapLayout struct {
	CellSize fyne.Size
	colCount int
	rowCount int
}

// Returns a new DynamicGridWrapLayout instance
// nearly identical to GridWrapLayout but where each child
// can be resized dynamically to fit content
func NewDynamicGridWrapLayout(size fyne.Size) fyne.Layout {
	return &dynamicGridWrapLayout{size, 1, 1}
}

// Layout is called to pack all child objects into a specified size.
// For a GridWrapLayout this will attempt to lay all the child objects in a row
// and wrap to a new row if the size is not large enough.
func (g *dynamicGridWrapLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	padding := theme.Padding()
	g.colCount = 1
	g.rowCount = 0

	if size.Width > g.CellSize.Width {
		g.colCount = int(math.Floor(float64(size.Width+padding) / float64(g.CellSize.Width+padding)))
	}

	i, x, y, maxHeight := 0, float32(0), float32(0), float32(0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if i%g.colCount == 0 {
			g.rowCount++
		}

		child.Move(fyne.NewPos(x, y))
		child.Resize(fyne.NewSize(g.CellSize.Width, child.MinSize().Height))
		// Store height of largest element in a row
		if child.MinSize().Height > maxHeight {
			maxHeight = child.MinSize().Height
		}

		if (i+1)%g.colCount == 0 {
			x = 0
			// Pad based on largest element in previous row
			y += maxHeight + padding
			maxHeight = 0
		} else {
			x += g.CellSize.Width + padding
		}
		i++
	}
}

func (g *dynamicGridWrapLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	rows := max(g.rowCount, 1)
	return fyne.NewSize(g.CellSize.Width, (g.CellSize.Height*float32(rows))+(float32(rows-1)*theme.Padding()))
}
