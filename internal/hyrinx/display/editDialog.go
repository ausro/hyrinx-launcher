package display

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AppDetails struct {
	Name  string
	Path  string
	Image string
	Opts  string
}

func EditAppDialog(parent fyne.Window, isEdit bool, initialValues *AppDetails, onSave func(*AppDetails)) {
	name := ""
	path := ""
	image := ""
	opts := ""

	if isEdit && initialValues != nil {
		name = initialValues.Name
		path = initialValues.Path
		image = initialValues.Image
		opts = initialValues.Opts
	}

	// Name field
	nameEntry := widget.NewEntry()
	nameEntry.SetText(name)
	nameField := widget.NewFormItem("Name", nameEntry)

	// Path field with file browser button
	pathEntry := widget.NewEntry()
	pathEntry.SetText(path)
	pathEntry.Validator = func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("path cannot be empty")
		}
		return nil
	}
	pathBtn := widget.NewButtonWithIcon("", theme.FileApplicationIcon(), func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				return
			}
			if reader == nil {
				return
			}
			_, err = storage.Exists(reader.URI())
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}
			pathEntry.SetText(reader.URI().Path())
			pathEntry.Refresh()
		}, parent)
	})
	pathBtn.Resize(fyne.NewSize(32, 32))
	pathField := widget.NewFormItem("Path", container.NewBorder(nil, nil, nil, pathBtn, pathEntry))

	// Image field with file browser button
	imageEntry := widget.NewEntry()
	imageEntry.SetText(image)
	imageBtn := widget.NewButtonWithIcon("", theme.FileImageIcon(), func() {
		f := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				return
			}
			if reader == nil {
				return
			}
			res, err := fyne.LoadResourceFromPath(reader.URI().Path())
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}
			_, err = storage.Exists(reader.URI())
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}
			imageEntry.SetText(reader.URI().Path())
			imageEntry.Icon = res
			imageEntry.Refresh()
		}, parent)
		f.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpeg", ".jpg"}))
		f.Show()
	})
	imageBtn.Resize(fyne.NewSize(32, 32))
	imageField := widget.NewFormItem("Image", container.NewBorder(nil, nil, nil, imageBtn, imageEntry))

	// Options field
	optsEntry := widget.NewEntry()
	optsEntry.SetText(opts)
	optsField := widget.NewFormItem("Launch args", optsEntry)

	items := []*widget.FormItem{nameField, pathField, imageField, optsField}

	title := "Edit App"
	if !isEdit {
		title = "Create App"
	}

	d := dialog.NewForm(title, "Save", "Cancel", items, func(ok bool) {
		if !ok {
			return
		}
		onSave(&AppDetails{
			Name:  nameEntry.Text,
			Path:  pathEntry.Text,
			Image: imageEntry.Text,
			Opts:  optsEntry.Text,
		})
	}, parent)

	d.Resize(fyne.NewSize(600, 300))
	d.Show()
}
