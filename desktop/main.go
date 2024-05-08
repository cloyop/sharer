package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MakeApp() fyne.Window {
	a := app.New()
	w := a.NewWindow("GoSharer")
	w.SetFixedSize(false)
	w.Resize(fyne.NewSize(500, 400))
	w.SetContent(ChoosePathGUI(w))
	return w
}
func ChoosePathGUI(w fyne.Window) fyne.CanvasObject {
	content := widget.NewLabel("Want to share or recieve")
	content.Alignment = fyne.TextAlignCenter
	content.TextStyle = widget.RichTextStyleCodeInline.TextStyle
	shareButton := widget.NewButton("Share", func() {
		w.SetContent(MakeShareGUI(newShareData(), w))
	})
	recieveButton := widget.NewButton("Recieve", func() {
		w.SetContent(MakeReceiveGUI(NewRecieveData(), w))
	})
	c := container.NewCenter(container.NewVBox(content, container.NewVBox(shareButton, recieveButton)))
	return c
}
