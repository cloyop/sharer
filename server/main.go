package server

import (
	"log"
	"net"

	pb "github.com/cloyop/sharer/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedShareServer
}

func RecieveHandler(args []string) {
	// size (optional) - default 20mb
	// port (optional) - default 9000
	// Process args
	port := "9000"
	sizeInt := 20 * 1_000_000
	strToken := makeSessionToken()
	OpenToReceive(strToken, port, sizeInt)
}

func OpenToReceive(token, port string, size int) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(grpc.MaxRecvMsgSize(size), grpc.MaxSendMsgSize(size))
	pb.RegisterShareServer(s, &Server{})
	log.Printf("server listening at port :%v", port)
	log.Printf("Clients will need authenticate with the following Token: \n  %v \n", token)
	if err := s.Serve(ln); err != nil {
		log.Fatal(err)
	}
}
