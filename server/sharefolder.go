package server

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"

	pb "github.com/cloyop/sharer/proto"
	"google.golang.org/grpc/peer"
)

func (s *Server) ShareFolder(ctx context.Context, in *pb.ShareFolderRequest) (*pb.ShareResponse, error) {
	if s.AuthRequired {
		if s.AuthToken != in.Token {
			return &pb.ShareResponse{Message: "Token Invalid"}, nil
		}
	}
	p, _ := peer.FromContext(ctx)
	folder := in.GetFolder()
	log.Printf("Share request from %v. Folder: %v of %d bytes\n", p.Addr, folder.Name, folder.Size)
	if err := unFold(folder); err != nil {
		return &pb.ShareResponse{Message: err.Error()}, nil
	}
	log.Printf("Folder recieve and save successfully %v of %d bytes \n", folder.Name, folder.Size)
	return &pb.ShareResponse{Message: fmt.Sprintf("File Shared %v", folder.Name), Success: true, BytesShared: int64(folder.Size)}, nil
}
func unFold(folder *pb.Folder) error {
	if err := mustCreateFolder(&folder.Name); err != nil {
		return err
	}
	fmt.Println(">", folder.Name)
	for _, file := range folder.Files {
		if err := unFile(file); err != nil {
			log.Printf("Cannot create: %v due %v", file.Name, err)
		}
	}
	for _, subFolder := range folder.Folders {
		if err := unFold(subFolder); err != nil {
			log.Printf("Cannot create: %v due %v", folder.Name, err)
		}
	}
	return nil
}
func mustCreateFolder(folderName *string) error {
	exist, err := fileIsExists(*folderName)
	if err != nil {
		return err
	}
	if exist {
		nn := fmt.Sprintf("%v%d", *folderName, rand.Intn(15))
		log.Printf("File %v exists, creating it as %v", *folderName, nn)
		*folderName = nn
	}
	if err = os.Mkdir(*folderName, os.ModePerm); err != nil {
		return err
	}
	return nil
}
