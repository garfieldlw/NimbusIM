// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: unique/unique.proto

package unique

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

const (
	Unique_GetUniqueId_FullMethodName     = "/com.NimbusIM.proto.unique.Unique/GetUniqueId"
	Unique_GetUniqueIdInfo_FullMethodName = "/com.NimbusIM.proto.unique.Unique/GetUniqueIdInfo"
)

// UniqueClient is the client API for Unique service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UniqueClient interface {
	GetUniqueId(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	GetUniqueIdInfo(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error)
}

type uniqueClient struct {
	cc grpc.ClientConnInterface
}

func NewUniqueClient(cc grpc.ClientConnInterface) UniqueClient {
	return &uniqueClient{cc}
}

func (c *uniqueClient) GetUniqueId(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Unique_GetUniqueId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uniqueClient) GetUniqueIdInfo(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error) {
	out := new(InfoResponse)
	err := c.cc.Invoke(ctx, Unique_GetUniqueIdInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UniqueServer is the server API for Unique service.
// All implementations should embed UnimplementedUniqueServer
// for forward compatibility
type UniqueServer interface {
	GetUniqueId(context.Context, *Request) (*Response, error)
	GetUniqueIdInfo(context.Context, *InfoRequest) (*InfoResponse, error)
}

// UnimplementedUniqueServer should be embedded to have forward compatible implementations.
type UnimplementedUniqueServer struct {
}

func (UnimplementedUniqueServer) GetUniqueId(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUniqueId not implemented")
}
func (UnimplementedUniqueServer) GetUniqueIdInfo(context.Context, *InfoRequest) (*InfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUniqueIdInfo not implemented")
}

// UnsafeUniqueServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UniqueServer will
// result in compilation errors.
type UnsafeUniqueServer interface {
	mustEmbedUnimplementedUniqueServer()
}

func RegisterUniqueServer(s grpc.ServiceRegistrar, srv UniqueServer) {
	s.RegisterService(&Unique_ServiceDesc, srv)
}

func _Unique_GetUniqueId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UniqueServer).GetUniqueId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Unique_GetUniqueId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UniqueServer).GetUniqueId(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Unique_GetUniqueIdInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UniqueServer).GetUniqueIdInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Unique_GetUniqueIdInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UniqueServer).GetUniqueIdInfo(ctx, req.(*InfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Unique_ServiceDesc is the grpc.ServiceDesc for Unique service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Unique_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.NimbusIM.proto.unique.Unique",
	HandlerType: (*UniqueServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUniqueId",
			Handler:    _Unique_GetUniqueId_Handler,
		},
		{
			MethodName: "GetUniqueIdInfo",
			Handler:    _Unique_GetUniqueIdInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "unique/unique.proto",
}
