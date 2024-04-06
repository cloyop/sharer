package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/cloyop/sharer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SharerHandler(args []string) {
	// need token to auth
	// address to dial
	// file | folder path
	// processArgs
	filePath := ""
	address := ""

	filePath = strings.TrimRight(filePath, "/")
	fInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	if fInfo.IsDir() {
		shareFolder(&pb.ShareFolderRequest{Folder: makeShareFolder(filePath)}, address)
	} else {
		shareFile(&pb.ShareFileRequest{File: makeShareFile(filePath)}, address)
	}
}
func mustDial(addr string) (*grpc.ClientConn, context.Context, context.CancelFunc) {
	// grpc auth
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not Dial: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)

	return conn, ctx, cancel
}
func readResponse(shr *pb.ShareResponse, err error) {
	if err != nil {
		log.Fatalf("Transport Error: %v\n", err.Error())
	}
	if !shr.Success {
		log.Fatalf("could not share successfully: %v", shr.Message)
	}
	fmt.Printf("Operation Sucess: %d bytes shared\n", shr.BytesShared)
}
