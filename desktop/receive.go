package desktop

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/cloyop/sharer/fs"
	"github.com/cloyop/sharer/proto"
)

type RecieveData struct {
	useGBSize,
	dropLocationSetted,
	insecureMode bool
	token,
	dropLocation string
	sizeEntry  *widget.Entry
	tokenEntry *widget.Entry
}

func NewRecieveData() *RecieveData {
	return &RecieveData{
		sizeEntry:  entryField("Max size per request ( Default 20MB )"),
		tokenEntry: entryField("Custom token ( Default Auto Generated )"),
	}
}
func MakeReceiveGUI(data *RecieveData, w fyne.Window) fyne.CanvasObject {
	checkGB := widget.NewCheck("Max size on GB measurement? ( Default on MB )", func(value bool) { data.useGBSize = value })

	tb := widget.NewToolbar(
		widget.NewToolbarAction(theme.NavigateBackIcon(), func() { w.SetContent(ChoosePathGUI(w)) }),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() { w.SetContent(MakeReceiveGUI(NewRecieveData(), w)) }),
	)
	size := container.NewVBox(container.NewCenter(checkGB), data.sizeEntry)
	customTokenMessage := "Used to authenticate incoming request,\nthe requests without it will be rejected.\nIt will be generated automatically.\ncan use your preferred custom token tho (min 16 length)."
	DropBox := container.NewVBox(selectDropLocation(data, w))
	if data.dropLocationSetted {
		DropBox.Add(centerLabel(data.dropLocation))
	}
	return container.NewCenter(
		container.NewVBox(
			tb,
			DropBox,
			size,
			container.NewBorder(base(), base(), centerLabel("Using a custom token"), infoButton(widget.NewLabel(customTokenMessage), w)),
			data.tokenEntry,
			centeredY(container.NewCenter(runButton(data, w))),
		),
	)
}

func selectDropLocation(data *RecieveData, w fyne.Window) *widget.Button {
	var onTFolder = func(uc fyne.ListableURI, err error) {
		if err != nil {
			genericModal(centerLabel(err.Error()), w, nil).Show()
			return
		}
		if uc != nil {
			data.dropLocation = uc.Path()
			data.dropLocationSetted = true
			w.SetContent(MakeReceiveGUI(data, w))
		}
	}
	folderButton := widget.NewButtonWithIcon("Select drop location", theme.FolderIcon(), func() {
		dialog.ShowFolderOpen(onTFolder, w)
	})
	return folderButton
}
func runButton(data *RecieveData, w fyne.Window) *widget.Button {
	run := func() {
		if data.tokenEntry.Text != "" {
			if len(data.tokenEntry.Text) < 16 {
				genericModal(centerLabel("Using custom token has to be min 16 character length"), w, nil).Show()
				return
			}
			data.token = data.tokenEntry.Text
		}

		if !data.dropLocationSetted {
			genericModal(centerLabel("Missing drop location"), w, nil).Show()
			return
		}
		var size int = 20 * int(fs.MB)
		if data.sizeEntry.Text != "" {
			n, err := strconv.Atoi(data.sizeEntry.Text)
			if err != nil {
				fmt.Println(err)
				return
			}
			if data.useGBSize {
				size = n * int(fs.GB)
			} else {
				size = n * int(fs.MB)
			}
		}
		srv := proto.RecieveSrv(size, data.token)
		m := widget.NewModalPopUp(centerLabel("Setting up server"), w.Canvas())
		m.Show()
		<-srv.Running
		m.Hide()
		c := container.NewVBox(
			centerLabel("Waiting on incoming requests"),
			centerLabel("Drop location : "+data.dropLocation),
			container.NewBorder(base(), base(), widget.NewLabel("Address"), copyToClipButton(srv.Addr, w)),
		)
		if !data.insecureMode {
			c.Add(container.NewBorder(base(), base(), widget.NewLabel("Auth Token"), copyToClipButton(srv.Token, w)))
		}
		c.Add(centeredX(widget.NewButton("Cancel", func() {
			srv.Close()
			m.Hide()
		})))
		m = widget.NewModalPopUp(container.NewCenter(
			c,
		),
			w.Canvas(),
		)
		m.Show()
		go func(data *RecieveData, srv *proto.Server, w fyne.Window, m *widget.PopUp) {
			for req := range srv.RequestChan {
				dirName := fmt.Sprintf("%v/goshare_%d/", data.dropLocation, time.Now().Unix())
				if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
					fmt.Println(err)
					return
				}
				if req.Type == proto.FileRequest && req.File != nil {
					fm := genericModal(
						centerLabel(fmt.Sprintf("Incomming request\nFile %v of %s", req.File.Name, fs.BytesToScalaBytes(len(req.File.Payload)))),
						w,
						nil)
					fm.Show()
					if err := req.File.UnFile(dirName); err != nil {
						fmt.Println(err)
						return
					}
					fm.Hide()
					genericModal(centerLabel(fmt.Sprintf("File %v received Succesfully\n%v", req.File.Name, dirName)), w, nil).Show()
				}
				if req.Type == proto.FolderRequest && req.Folder != nil {
					fm := genericModal(
						centerLabel(fmt.Sprintf("Incomming request\nFolder %v of %s", req.Folder.Name, fs.BytesToScalaBytes(int(req.Folder.Size)))),
						w,
						nil)
					fm.Show()
					if err := req.Folder.UnFold(dirName); err != nil {
						fmt.Println(err)
						return
					}
					fm.Hide()
					genericModal(centerLabel(fmt.Sprintf("Folder %v received Succesfully\n%v", req.Folder.Name, dirName)), w, nil).Show()
				}
			}
		}(data, srv, w, m)

	}
	return widget.NewButton("Run", run)
}
