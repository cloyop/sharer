// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/file.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ShareClient is the client API for Share service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShareClient interface {
	Share(ctx context.Context, in *ShareRequest, opts ...grpc.CallOption) (*ShareResponse, error)
}

type shareClient struct {
	cc grpc.ClientConnInterface
}

func NewShareClient(cc grpc.ClientConnInterface) ShareClient {
	return &shareClient{cc}
}

func (c *shareClient) Share(ctx context.Context, in *ShareRequest, opts ...grpc.CallOption) (*ShareResponse, error) {
	out := new(ShareResponse)
	err := c.cc.Invoke(ctx, "/proto.Share/Share", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShareServer is the server API for Share service.
// All implementations must embed UnimplementedShareServer
// for forward compatibility
type ShareServer interface {
	Share(context.Context, *ShareRequest) (*ShareResponse, error)
	mustEmbedUnimplementedShareServer()
}

// UnimplementedShareServer must be embedded to have forward compatible implementations.
type UnimplementedShareServer struct {
}

func (UnimplementedShareServer) Share(context.Context, *ShareRequest) (*ShareResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Share not implemented")
}
func (UnimplementedShareServer) mustEmbedUnimplementedShareServer() {}

// UnsafeShareServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShareServer will
// result in compilation errors.
type UnsafeShareServer interface {
	mustEmbedUnimplementedShareServer()
}

func RegisterShareServer(s grpc.ServiceRegistrar, srv ShareServer) {
	s.RegisterService(&Share_ServiceDesc, srv)
}

func _Share_Share_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShareRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShareServer).Share(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Share/Share",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShareServer).Share(ctx, req.(*ShareRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Share_ServiceDesc is the grpc.ServiceDesc for Share service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Share_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Share",
	HandlerType: (*ShareServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Share",
			Handler:    _Share_Share_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/file.proto",
}
