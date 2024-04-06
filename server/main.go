package server

import (
	"fmt"
	"log"
	"net"

	pb "github.com/cloyop/sharer/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedShareServer
	AuthRequired bool
	AuthToken    string
}

func RecieveHandler(size int, port, token string, noAuth bool) {
	srv := &Server{}
	if !noAuth {
		srv.AuthRequired = true
		srv.AuthToken = token
	}
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(grpc.MaxRecvMsgSize(size), grpc.MaxSendMsgSize(size))
	pb.RegisterShareServer(s, srv)
	fmt.Printf("server running at port :%v", port)
	if !noAuth {
		fmt.Printf("Clients will need authenticate with the following Token: %v\n", token)
	}
	if err := s.Serve(ln); err != nil {
		log.Fatal(err)
	}
}
