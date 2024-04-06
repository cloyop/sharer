package client

import (
	"fmt"
	"os"

	pb "github.com/cloyop/sharer/proto"
)

func shareFolder(r *pb.ShareFolderRequest, address string) {
	conn, ctx, cancel := mustDial(address)
	defer cancel()
	shareClient := pb.NewShareClient(conn)
	readResponse(shareClient.ShareFolder(ctx, r))
}

func makeShareFolder(folderName string) *pb.Folder {
	files, err := os.ReadDir(folderName)
	if err != nil {
		fmt.Printf("Could not read Dir: %v due: %v", folderName, err)
		return nil
	}
	var folder pb.Folder
	folder.Name = folderName
	for _, f := range files {
		if f.IsDir() {
			if shares := makeShareFolder(folderName + "/" + f.Name()); shares != nil {
				folder.Folders = append(folder.Folders, shares)
				folder.Size += shares.Size
			}
		} else {
			if shareFile := makeShareFile(folderName + "/" + f.Name()); shareFile != nil {
				folder.Files = append(folder.Files, shareFile)
				folder.Size += uint64(len(shareFile.Payload))
			}
		}
	}
	return &folder
}
