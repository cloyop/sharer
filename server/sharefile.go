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

func (s *Server) ShareFile(ctx context.Context, in *pb.ShareFileRequest) (*pb.ShareResponse, error) {
	if s.AuthRequired {
		if s.AuthToken != in.Token {
			return &pb.ShareResponse{Message: "Token Invalid"}, nil
		}
	}
	p, _ := peer.FromContext(ctx)
	file := in.GetFile()
	log.Printf("Share request from %v. File: %v of %d bytes\n", p.Addr, file.Name, len(file.Payload))
	if err := unFile(file); err != nil {
		return &pb.ShareResponse{Message: err.Error()}, nil
	}
	log.Printf("File recieve and save successfully %v of %d bytes \n", file.Name, int64(len(file.Payload)))
	return &pb.ShareResponse{Message: fmt.Sprintf("File Shared %v", file.Name), Success: true, BytesShared: int64(len(file.Payload))}, nil
}
func unFile(file *pb.File) error {
	f, err := mustCreateFile(&file.Name)
	if err != nil {
		return err
	}
	fmt.Println(" -", file.Name)
	defer f.Close()
	_, err = f.Write(file.Payload)
	if err != nil {
		return err
	}
	return nil
}
func mustCreateFile(fn *string) (*os.File, error) {
	exist, err := fileIsExists(*fn)
	if err != nil {
		return nil, err
	}
	if exist {
		nn := fmt.Sprintf("%d%v", rand.Intn(15), *fn)
		fmt.Printf("File %v exists, creating it as %v", *fn, nn)
		*fn = nn
	}
	f, err := os.Create(*fn)
	if err != nil {
		return nil, err
	}
	return f, nil
}
