package fs

import (
	"os"
	"strings"
)

type File struct {
	Name    string
	Path    string
	Payload []byte
	Ext     string
}

func makeFile(name, path, location string) (*File, error) {
	data, err := os.ReadFile(location + path + name)
	if err != nil {
		return nil, err
	}

	return &File{
		Name:    name,
		Path:    path,
		Payload: data,
		Ext:     getExt(name),
	}, nil
}
func (file *File) UnFile(location string) error {
	f, err := CreateFile(location + file.Path + file.Name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(file.Payload)
	if err != nil {
		return err
	}
	return nil
}

func getExt(name string) (ext string) {
	if strings.Contains(name, ".") {
		slic := strings.Split(name, ".")
		ext = "." + slic[len(slic)-1]
	}
	return
}
func CreateFile(fn string) (*os.File, error) {
	f, err := os.Create(fn)
	if err != nil {
		return nil, err
	}
	return f, nil
}
