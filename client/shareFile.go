package client

import (
	"fmt"
	"os"
	"strings"

	pb "github.com/cloyop/sharer/proto"
)

func shareFile(r *pb.ShareFileRequest, address string) {
	conn, ctx, cancel := mustDial(address)
	defer cancel()
	shareClient := pb.NewShareClient(conn)
	readResponse(shareClient.ShareFile(ctx, r))
}
func makeShareFile(fileName string) *pb.File {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Couldnt Ship file: %v due: %v", fileName, err)
		return nil
	}
	filetype := ""
	if strings.Contains(fileName, ".") {
		slic := strings.Split(fileName, ".")
		filetype = "." + slic[len(slic)-1]
	}
	return &pb.File{
		Name:    fileName,
		Type:    filetype,
		Payload: bytes,
	}
}
