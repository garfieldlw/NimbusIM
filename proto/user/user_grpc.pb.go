// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: user/user.proto

package user

import (
	common "github.com/garfieldlw/NimbusIM/proto/common"
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
	User_Upset_FullMethodName        = "/com.NimbusIM.proto.user.User/Upset"
	User_UpdateStatus_FullMethodName = "/com.NimbusIM.proto.user.User/UpdateStatus"
	User_DetailById_FullMethodName   = "/com.NimbusIM.proto.user.User/DetailById"
)

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserClient interface {
	Upset(ctx context.Context, in *UserUpsetRequest, opts ...grpc.CallOption) (*UserEntity, error)
	UpdateStatus(ctx context.Context, in *UpdateStatusRequest, opts ...grpc.CallOption) (*common.Empty, error)
	DetailById(ctx context.Context, in *common.DetailByIdRequest, opts ...grpc.CallOption) (*UserEntity, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) Upset(ctx context.Context, in *UserUpsetRequest, opts ...grpc.CallOption) (*UserEntity, error) {
	out := new(UserEntity)
	err := c.cc.Invoke(ctx, User_Upset_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) UpdateStatus(ctx context.Context, in *UpdateStatusRequest, opts ...grpc.CallOption) (*common.Empty, error) {
	out := new(common.Empty)
	err := c.cc.Invoke(ctx, User_UpdateStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) DetailById(ctx context.Context, in *common.DetailByIdRequest, opts ...grpc.CallOption) (*UserEntity, error) {
	out := new(UserEntity)
	err := c.cc.Invoke(ctx, User_DetailById_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServer is the server API for User service.
// All implementations should embed UnimplementedUserServer
// for forward compatibility
type UserServer interface {
	Upset(context.Context, *UserUpsetRequest) (*UserEntity, error)
	UpdateStatus(context.Context, *UpdateStatusRequest) (*common.Empty, error)
	DetailById(context.Context, *common.DetailByIdRequest) (*UserEntity, error)
}

// UnimplementedUserServer should be embedded to have forward compatible implementations.
type UnimplementedUserServer struct {
}

func (UnimplementedUserServer) Upset(context.Context, *UserUpsetRequest) (*UserEntity, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Upset not implemented")
}
func (UnimplementedUserServer) UpdateStatus(context.Context, *UpdateStatusRequest) (*common.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStatus not implemented")
}
func (UnimplementedUserServer) DetailById(context.Context, *common.DetailByIdRequest) (*UserEntity, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetailById not implemented")
}

// UnsafeUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServer will
// result in compilation errors.
type UnsafeUserServer interface {
	mustEmbedUnimplementedUserServer()
}

func RegisterUserServer(s grpc.ServiceRegistrar, srv UserServer) {
	s.RegisterService(&User_ServiceDesc, srv)
}

func _User_Upset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserUpsetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Upset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_Upset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Upset(ctx, req.(*UserUpsetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_UpdateStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).UpdateStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_UpdateStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).UpdateStatus(ctx, req.(*UpdateStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_DetailById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.DetailByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).DetailById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_DetailById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).DetailById(ctx, req.(*common.DetailByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// User_ServiceDesc is the grpc.ServiceDesc for User service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var User_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.NimbusIM.proto.user.User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Upset",
			Handler:    _User_Upset_Handler,
		},
		{
			MethodName: "UpdateStatus",
			Handler:    _User_UpdateStatus_Handler,
		},
		{
			MethodName: "DetailById",
			Handler:    _User_DetailById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user/user.proto",
}
