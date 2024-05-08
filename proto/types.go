package proto

import (
	"bytes"
	"encoding/gob"
	"net"

	"github.com/cloyop/sharer/fs"
	"github.com/cloyop/sharer/utils"
)

var FileRequest uint8 = 1
var FolderRequest uint8 = 2

type Server struct {
	UnimplementedShareServer
	Token       string
	RequestChan chan *Request
	Running     chan bool
	ln          net.Listener
	port        string
	size        int
	Addr        string
}
type Request struct {
	Addr   net.Addr
	Type   uint8
	File   *fs.File
	Folder *fs.Folder
}

func newRequest() *Request {
	return &Request{
		File:   &fs.File{},
		Folder: &fs.Folder{Files: []*fs.File{}, Folders: []*fs.Folder{}},
	}
}
func (r *Request) Encrypt(token string) ([]byte, error) {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(r); err != nil {
		return nil, err
	}
	return utils.Cipher([]byte(token), b.Bytes())
}
func (r *Request) DesEncrypt(data []byte, token string) error {
	data, err := utils.UnCipher([]byte(token), data)
	if err != nil {
		return err
	}
	b := new(bytes.Buffer)
	b.Write(data)
	if err := gob.NewDecoder(b).Decode(r); err != nil {
		return err
	}
	return nil
}
