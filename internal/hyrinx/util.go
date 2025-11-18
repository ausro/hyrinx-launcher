package hyrinx

import (
	"strings"

	"fyne.io/fyne/v2"
)

// Checks if the given item is a windows executable
//   - path: The absolute path to the file
func isWindowsExecutable(path string) bool {
	return strings.HasSuffix(path, ".exe")
}

func GetRootWindow() *fyne.Window {
	// Get first window
	if w := fyne.CurrentApp().Driver().AllWindows(); len(w) > 0 {
		return &w[0]
	}

	return nil
}
