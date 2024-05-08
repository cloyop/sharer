package desktop

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/cloyop/sharer/fs"
	"github.com/cloyop/sharer/proto"
)

type SharingData struct {
	ShareFile    *fs.File
	ShareFolder  *fs.Folder
	AddressEntry *widget.Entry
	TokenEntry   *widget.Entry
	Loaded,
	FolderSelected,
	FileSelected bool
}

func (sd *SharingData) resetPayload() {
	sd.ShareFile = &fs.File{}
	sd.ShareFolder = &fs.Folder{}
	sd.FileSelected = false
	sd.FolderSelected = false
	sd.Loaded = false
}
func newShareData() *SharingData {
	return &SharingData{
		ShareFile:    &fs.File{},
		ShareFolder:  &fs.Folder{},
		AddressEntry: entryField("Address"),
		TokenEntry:   entryField("Token"),
	}
}

func MakeShareGUI(data *SharingData, w fyne.Window) fyne.CanvasObject {
	tb := widget.NewToolbar(
		widget.NewToolbarAction(theme.NavigateBackIcon(), func() { w.SetContent(ChoosePathGUI(w)) }),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() { w.SetContent(MakeShareGUI(newShareData(), w)) }),
	)
	var label fyne.CanvasObject = base()
	if data.Loaded {
		if data.FileSelected {
			label = centerLabel(fmt.Sprintf("File %v loaded %v", data.ShareFile.Name, fs.BytesToScalaBytes(len(data.ShareFile.Payload))))
		} else {
			label = centerLabel(fmt.Sprintf("Folder %v loaded %v", data.ShareFolder.Name, fs.BytesToScalaBytes(int(data.ShareFolder.Size))))
		}
	}
	var top = container.NewBorder(tb, label, (fileButton(data, w)), (folderButton(data, w)))
	var middle = container.NewVBox(centerLabel("Address to send: "), data.AddressEntry, centerLabel("Token to authenticate with receiver: "), data.TokenEntry)
	c := container.NewCenter(container.NewVBox(
		top,
		middle,
		centeredY(container.NewCenter(sendButton(data, w))),
	),
	)
	return c
}

func sendButton(data *SharingData, w fyne.Window) *widget.Button {
	var onTSend = func() {
		if data.AddressEntry.Text == "" {
			genericModal(centerLabel("Missing Address"), w, nil).Show()
			return
		}
		if data.TokenEntry.Text == "" {
			genericModal(centerLabel("Missing Token"), w, nil).Show()
			return
		}
		if !data.Loaded {
			genericModal(centerLabel("Missing data to share"), w, nil).Show()
			return
		}
		r := &proto.Request{}
		if data.FileSelected {
			if len(data.ShareFile.Payload) == 0 {
				genericModal(centerLabel("Missing file data to share"), w, nil).Show()
				return
			}
			r.File = data.ShareFile
			r.Type = proto.FileRequest
		} else if data.FolderSelected {
			if data.ShareFolder.Size == 0 {
				genericModal(centerLabel("Missing folder data to share"), w, nil).Show()
				return
			}
			r.Folder = data.ShareFolder
			r.Type = proto.FolderRequest
		} else {
			genericModal(centerLabel("Missing data to share"), w, nil).Show()
			return
		}
		req, errReq := r.Encrypt(data.TokenEntry.Text)
		if errReq != nil {
			fmt.Println(errReq)
		}
		response, err := proto.Share(&proto.ShareRequest{Data: req}, data.AddressEntry.Text)
		if err != nil {
			genericModal(centerLabel("Internal Error "+err.Error()), w, nil).Show()
			return
		}
		if !response.Success {
			genericModal(centerLabel("Error: "+response.Message), w, nil).Show()
			return
		}
		data.resetPayload()
		genericModal(centerLabel(response.Message), w, func() { w.SetContent(MakeShareGUI(data, w)) }).Show()
	}
	return widget.NewButtonWithIcon("", theme.MailSendIcon(), onTSend)
}
func folderButton(data *SharingData, w fyne.Window) *widget.Button {
	var onTFolder = func(uc fyne.ListableURI, err error) {
		if data.FileSelected {
			data.ShareFile = &fs.File{}
			data.FileSelected = false
			data.Loaded = false
		}
		if err != nil {
			genericModal(centerLabel(err.Error()), w, nil).Show()
			return
		}
		if uc != nil {
			if uc.Path() == "/" {
				genericModal(centerLabel("Cannot read this directory"), w, nil).Show()
				return
			}
			if uc.Path() == os.Getenv("HOME") || uc.Path() == "/home" {
				go func() {
					l := centerLabel("You are about to load a entire delicated directory:\n" + uc.Path() + "\n Want to proceed?")
					y := widget.NewButton("Yes", nil)
					n := widget.NewButton("No", nil)
					pop := widget.NewModalPopUp(container.NewVBox(l, container.NewCenter(container.NewHBox(y, n))), w.Canvas())
					y.OnTapped = func() {
						pop.Hide()
						LoadFold(data, w, uc.Path(), uc.Name())
					}
					n.OnTapped = pop.Hide
					pop.Show()
				}()
				return
			}
			LoadFold(data, w, uc.Path(), uc.Name())
		}
	}
	folderButton := widget.NewButtonWithIcon("Select a Folder", theme.FolderIcon(), func() {
		dialog.ShowFolderOpen(onTFolder, w)
	})
	return folderButton
}
func fileButton(data *SharingData, w fyne.Window) *widget.Button {
	var onTFile = func(uc fyne.URIReadCloser, err error) {
		if data.FolderSelected {
			data.ShareFolder = &fs.Folder{}
			data.FolderSelected = false
			data.Loaded = false
		}
		if err != nil {
			genericModal(centerLabel(err.Error()), w, nil).Show()
			return
		}
		if uc != nil {
			label := centerLabel(fmt.Sprintf("Loading file %v", uc.URI().Name()))
			pop := widget.NewModalPopUp(container.NewCenter(label), w.Canvas())
			pop.Show()
			defer pop.Hide()
			b := new(bytes.Buffer)
			wr, err := io.Copy(b, uc)
			if err != nil {
				genericModal(centerLabel(err.Error()), w, nil).Show()
				return
			}
			if wr == 0 {
				genericModal(centerLabel(fmt.Sprintf("file ' %v ' is empty", uc.URI().Name())), w, nil).Show()
				return
			}
			if wr > int64(fs.GB*50) {
				genericModal(centerLabel(fmt.Sprintf("file ' %v ' too big", uc.URI().Name())), w, nil).Show()
				return
			}
			data.ShareFile.Name = uc.URI().Name()
			data.ShareFile.Ext = uc.URI().Extension()
			data.ShareFile.Payload = b.Bytes()
			data.FileSelected = true
			data.Loaded = true
			w.SetContent(MakeShareGUI(data, w))
		}
	}
	fileButton := widget.NewButtonWithIcon("Select a File", theme.FileIcon(), func() {
		dialog.ShowFileOpen(onTFile, w)
	})
	return fileButton
}
func LoadFold(data *SharingData, w fyne.Window, path, name string) {
	label := centerLabel(fmt.Sprintf("Loading Folders & Files from:\n%v", name))
	pop := widget.NewModalPopUp(container.NewCenter(label), w.Canvas())
	pop.Show()
	pat := strings.TrimRight(path, name)
	fold, err := fs.MakeShareFolder(name, pat)
	if err != nil {
		genericModal(centerLabel(err.Error()), w, nil).Show()
		return
	}
	*data.ShareFolder = *fold
	data.FolderSelected = true
	data.Loaded = true
	pop.Hide()
	genericModal(
		centerLabel(fmt.Sprintf("Folder %v loaded %v ", name, fs.BytesToScalaBytes(int(data.ShareFolder.Size)))),
		w,
		nil,
	).Show()
	w.SetContent(MakeShareGUI(data, w))

}
