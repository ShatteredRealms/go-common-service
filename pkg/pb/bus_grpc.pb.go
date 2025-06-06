// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.1
// source: sro/bus.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	BusService_ResetReaderBus_FullMethodName = "/sro.BusService/ResetReaderBus"
	BusService_ResetWriterBus_FullMethodName = "/sro.BusService/ResetWriterBus"
)

// BusServiceClient is the client API for BusService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BusServiceClient interface {
	ResetReaderBus(ctx context.Context, in *BusTarget, opts ...grpc.CallOption) (*ResetBusResponse, error)
	ResetWriterBus(ctx context.Context, in *BusTarget, opts ...grpc.CallOption) (*ResetBusResponse, error)
}

type busServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBusServiceClient(cc grpc.ClientConnInterface) BusServiceClient {
	return &busServiceClient{cc}
}

func (c *busServiceClient) ResetReaderBus(ctx context.Context, in *BusTarget, opts ...grpc.CallOption) (*ResetBusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResetBusResponse)
	err := c.cc.Invoke(ctx, BusService_ResetReaderBus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *busServiceClient) ResetWriterBus(ctx context.Context, in *BusTarget, opts ...grpc.CallOption) (*ResetBusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResetBusResponse)
	err := c.cc.Invoke(ctx, BusService_ResetWriterBus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BusServiceServer is the server API for BusService service.
// All implementations must embed UnimplementedBusServiceServer
// for forward compatibility.
type BusServiceServer interface {
	ResetReaderBus(context.Context, *BusTarget) (*ResetBusResponse, error)
	ResetWriterBus(context.Context, *BusTarget) (*ResetBusResponse, error)
	mustEmbedUnimplementedBusServiceServer()
}

// UnimplementedBusServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBusServiceServer struct{}

func (UnimplementedBusServiceServer) ResetReaderBus(context.Context, *BusTarget) (*ResetBusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetReaderBus not implemented")
}
func (UnimplementedBusServiceServer) ResetWriterBus(context.Context, *BusTarget) (*ResetBusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetWriterBus not implemented")
}
func (UnimplementedBusServiceServer) mustEmbedUnimplementedBusServiceServer() {}
func (UnimplementedBusServiceServer) testEmbeddedByValue()                    {}

// UnsafeBusServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BusServiceServer will
// result in compilation errors.
type UnsafeBusServiceServer interface {
	mustEmbedUnimplementedBusServiceServer()
}

func RegisterBusServiceServer(s grpc.ServiceRegistrar, srv BusServiceServer) {
	// If the following call pancis, it indicates UnimplementedBusServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BusService_ServiceDesc, srv)
}

func _BusService_ResetReaderBus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BusTarget)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusServiceServer).ResetReaderBus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BusService_ResetReaderBus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusServiceServer).ResetReaderBus(ctx, req.(*BusTarget))
	}
	return interceptor(ctx, in, info, handler)
}

func _BusService_ResetWriterBus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BusTarget)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BusServiceServer).ResetWriterBus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BusService_ResetWriterBus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BusServiceServer).ResetWriterBus(ctx, req.(*BusTarget))
	}
	return interceptor(ctx, in, info, handler)
}

// BusService_ServiceDesc is the grpc.ServiceDesc for BusService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BusService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sro.BusService",
	HandlerType: (*BusServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResetReaderBus",
			Handler:    _BusService_ResetReaderBus_Handler,
		},
		{
			MethodName: "ResetWriterBus",
			Handler:    _BusService_ResetWriterBus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sro/bus.proto",
}
