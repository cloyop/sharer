package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func base() *widget.BaseWidget {
	return &widget.BaseWidget{}
}
func centeredX(c fyne.CanvasObject) *fyne.Container {
	return container.NewGridWithColumns(3, base(), c, base())
}
func centeredY(c fyne.CanvasObject) *fyne.Container {
	return container.NewGridWithRows(3, base(), c, base())
}

func centerLabel(s string) *widget.Label {
	label := widget.NewLabel(s)
	label.Alignment = fyne.TextAlignCenter
	return label
}
func entryField(placeHolder string) *widget.Entry {
	e := widget.NewEntry()
	e.PlaceHolder = placeHolder
	return e
}
func genericModal(content *widget.Label, w fyne.Window, f func()) *widget.PopUp {
	popup := widget.NewModalPopUp(nil, w.Canvas())
	top := container.NewBorder(base(), base(), base(), widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		popup.Hide()
		if f != nil {
			f()
		}
	}))
	popup.Content = container.NewVBox(top, content)
	return popup
}
func copyToClipButton(s string, w fyne.Window) *widget.Button {
	return widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() { w.Clipboard().SetContent(s) })
}
func infoButton(content *widget.Label, w fyne.Window) *widget.Button {
	return widget.NewButtonWithIcon("", theme.InfoIcon(), func() { genericModal(content, w, nil).Show() })
}
