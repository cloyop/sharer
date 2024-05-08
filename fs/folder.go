package fs

import (
	"fmt"
	fsys "io/fs"
	"os"
	"sync"
)

type Folder struct {
	wg        *sync.WaitGroup
	foldersMu *sync.Mutex
	Folders   []*Folder
	Files     []*File
	Path      string
	Size      uint64
	Name      string
}

func (f *Folder) UnFold(location string) error {
	if err := os.Mkdir(location+f.Path, os.ModePerm); err != nil {
		return err
	}
	for _, v := range f.Folders {
		if err := v.UnFold(location); err != nil {
			fmt.Println(err)
			continue
		}
	}
	for _, v := range f.Files {
		if err := v.UnFile(location); err != nil {
			fmt.Println(err)
			continue
		}
	}
	return nil
}

func MakeShareFolder(name, location string) (*Folder, error) {
	fsStat, err := os.Stat(location + name)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if fsStat.Size() > int64(GB*50) {
		return nil, fmt.Errorf("size too large")
	}
	fs, err := os.ReadDir(location + name)
	if err != nil {
		return nil, err
	}
	folder := newFolder(name, name+"/")
	for _, f := range fs {
		if f.IsDir() {
			folder.wg.Add(1)
			go folder.appendFolder(location, f)
		} else {
			file, err := makeFile(f.Name(), folder.Path, location)
			if err != nil {
				fmt.Println(err)
				continue
			}
			folder.Files = append(folder.Files, file)
			folder.Size += uint64(len(file.Payload))
		}
	}
	folder.wg.Wait()
	return folder, nil
}
func makeFolder(name, path, location string) (*Folder, error) {
	fs, err := os.ReadDir(location + path + name)
	if err != nil {
		return nil, err
	}
	folder := newFolder(name, path+name+"/")
	for _, f := range fs {
		if f.IsDir() {
			folder.wg.Add(1)
			go folder.appendFolder(location, f)
		} else {
			file, err := makeFile(f.Name(), folder.Path, location)
			if err != nil {
				fmt.Println(err)
				continue
			}
			folder.Files = append(folder.Files, file)
			folder.Size += uint64(len(file.Payload))
		}
	}
	folder.wg.Wait()
	return folder, nil
}
func newFolder(name, path string) *Folder {
	return &Folder{
		wg:        &sync.WaitGroup{},
		foldersMu: &sync.Mutex{},
		Folders:   []*Folder{},
		Files:     []*File{},
		Name:      name,
		Path:      path,
	}
}
func (folder *Folder) appendFolder(location string, f fsys.DirEntry) {
	defer folder.wg.Done()
	fold, err := makeFolder(f.Name(), folder.Path, location)
	if err != nil {
		fmt.Println(err)
		return
	}
	folder.foldersMu.Lock()
	defer folder.foldersMu.Unlock()
	folder.Folders = append(folder.Folders, fold)
	folder.Size += fold.Size
}
