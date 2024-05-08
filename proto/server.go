package proto

import (
	context "context"
	"fmt"
	"log"
	"net"

	"github.com/cloyop/sharer/utils"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func (s *Server) Share(ctx context.Context, in *ShareRequest) (*ShareResponse, error) {
	data := in.GetData()
	r := newRequest()
	if err := r.DesEncrypt(data, s.Token); err != nil {
		return &ShareResponse{Message: "Invalid Credentials"}, nil
	}
	p, _ := peer.FromContext(ctx)
	r.Addr = p.Addr
	s.RequestChan <- r
	return &ShareResponse{Message: "Data Shared Successfully", Success: true}, nil
}
func RecieveSrv(size int, token string) *Server {
	srv := &Server{
		Token:       token,
		port:        ":9000",
		size:        size,
		RequestChan: make(chan *Request),
		Running:     make(chan bool),
	}
	ip, err := utils.GetIpAddr()
	if err != nil {
		fmt.Println(err)
	}
	srv.Addr = ip + srv.port
	if token == "" {
		srv.Token = utils.MakeSessionToken()
	}
	go srv.Run()
	return srv
}
func (srv *Server) Run() {
	var err error
	srv.ln, err = net.Listen("tcp", srv.port)

	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(grpc.MaxRecvMsgSize(srv.size), grpc.MaxSendMsgSize(srv.size))
	RegisterShareServer(s, srv)
	srv.Running <- true
	if err := s.Serve(srv.ln); err != nil {
		fmt.Println(err)
	}
}
func (srv *Server) Close() {
	srv.ln.Close()
}
