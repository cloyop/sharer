package proto

import (
	context "context"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func mustDial(addr string) (*grpc.ClientConn, context.Context, context.CancelFunc, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)

	return conn, ctx, cancel, nil
}
func Share(r *ShareRequest, address string) (*ShareResponse, error) {
	conn, ctx, cancel, err := mustDial(address)
	if err != nil {
		return nil, err
	}
	defer cancel()
	shareClient := NewShareClient(conn)
	return shareClient.Share(ctx, r)
}
