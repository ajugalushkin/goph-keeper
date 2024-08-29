// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: keeper/v1/keeper.proto

package v1

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
	KeeperServiceV1_RegisterV1_FullMethodName = "/keeper.v1.KeeperServiceV1/RegisterV1"
	KeeperServiceV1_LoginV1_FullMethodName    = "/keeper.v1.KeeperServiceV1/LoginV1"
	KeeperServiceV1_ListItems_FullMethodName  = "/keeper.v1.KeeperServiceV1/ListItems"
	KeeperServiceV1_SetItem_FullMethodName    = "/keeper.v1.KeeperServiceV1/SetItem"
)

// KeeperServiceV1Client is the client API for KeeperServiceV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeeperServiceV1Client interface {
	RegisterV1(ctx context.Context, in *RegisterRequestV1, opts ...grpc.CallOption) (*RegisterResponseV1, error)
	LoginV1(ctx context.Context, in *LoginRequestV1, opts ...grpc.CallOption) (*LoginResponseV1, error)
	ListItems(ctx context.Context, in *ListItemsRequest, opts ...grpc.CallOption) (*ListItemsResponse, error)
	SetItem(ctx context.Context, in *SetItemRequest, opts ...grpc.CallOption) (*SetItemResponse, error)
}

type keeperServiceV1Client struct {
	cc grpc.ClientConnInterface
}

func NewKeeperServiceV1Client(cc grpc.ClientConnInterface) KeeperServiceV1Client {
	return &keeperServiceV1Client{cc}
}

func (c *keeperServiceV1Client) RegisterV1(ctx context.Context, in *RegisterRequestV1, opts ...grpc.CallOption) (*RegisterResponseV1, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponseV1)
	err := c.cc.Invoke(ctx, KeeperServiceV1_RegisterV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperServiceV1Client) LoginV1(ctx context.Context, in *LoginRequestV1, opts ...grpc.CallOption) (*LoginResponseV1, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponseV1)
	err := c.cc.Invoke(ctx, KeeperServiceV1_LoginV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperServiceV1Client) ListItems(ctx context.Context, in *ListItemsRequest, opts ...grpc.CallOption) (*ListItemsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListItemsResponse)
	err := c.cc.Invoke(ctx, KeeperServiceV1_ListItems_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperServiceV1Client) SetItem(ctx context.Context, in *SetItemRequest, opts ...grpc.CallOption) (*SetItemResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetItemResponse)
	err := c.cc.Invoke(ctx, KeeperServiceV1_SetItem_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeeperServiceV1Server is the server API for KeeperServiceV1 service.
// All implementations must embed UnimplementedKeeperServiceV1Server
// for forward compatibility.
type KeeperServiceV1Server interface {
	RegisterV1(context.Context, *RegisterRequestV1) (*RegisterResponseV1, error)
	LoginV1(context.Context, *LoginRequestV1) (*LoginResponseV1, error)
	ListItems(context.Context, *ListItemsRequest) (*ListItemsResponse, error)
	SetItem(context.Context, *SetItemRequest) (*SetItemResponse, error)
	mustEmbedUnimplementedKeeperServiceV1Server()
}

// UnimplementedKeeperServiceV1Server must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedKeeperServiceV1Server struct{}

func (UnimplementedKeeperServiceV1Server) RegisterV1(context.Context, *RegisterRequestV1) (*RegisterResponseV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterV1 not implemented")
}
func (UnimplementedKeeperServiceV1Server) LoginV1(context.Context, *LoginRequestV1) (*LoginResponseV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginV1 not implemented")
}
func (UnimplementedKeeperServiceV1Server) ListItems(context.Context, *ListItemsRequest) (*ListItemsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListItems not implemented")
}
func (UnimplementedKeeperServiceV1Server) SetItem(context.Context, *SetItemRequest) (*SetItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetItem not implemented")
}
func (UnimplementedKeeperServiceV1Server) mustEmbedUnimplementedKeeperServiceV1Server() {}
func (UnimplementedKeeperServiceV1Server) testEmbeddedByValue()                         {}

// UnsafeKeeperServiceV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeeperServiceV1Server will
// result in compilation errors.
type UnsafeKeeperServiceV1Server interface {
	mustEmbedUnimplementedKeeperServiceV1Server()
}

func RegisterKeeperServiceV1Server(s grpc.ServiceRegistrar, srv KeeperServiceV1Server) {
	// If the following call pancis, it indicates UnimplementedKeeperServiceV1Server was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&KeeperServiceV1_ServiceDesc, srv)
}

func _KeeperServiceV1_RegisterV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServiceV1Server).RegisterV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeeperServiceV1_RegisterV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServiceV1Server).RegisterV1(ctx, req.(*RegisterRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeeperServiceV1_LoginV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServiceV1Server).LoginV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeeperServiceV1_LoginV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServiceV1Server).LoginV1(ctx, req.(*LoginRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeeperServiceV1_ListItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListItemsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServiceV1Server).ListItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeeperServiceV1_ListItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServiceV1Server).ListItems(ctx, req.(*ListItemsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeeperServiceV1_SetItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetItemRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServiceV1Server).SetItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeeperServiceV1_SetItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServiceV1Server).SetItem(ctx, req.(*SetItemRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KeeperServiceV1_ServiceDesc is the grpc.ServiceDesc for KeeperServiceV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KeeperServiceV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "keeper.v1.KeeperServiceV1",
	HandlerType: (*KeeperServiceV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterV1",
			Handler:    _KeeperServiceV1_RegisterV1_Handler,
		},
		{
			MethodName: "LoginV1",
			Handler:    _KeeperServiceV1_LoginV1_Handler,
		},
		{
			MethodName: "ListItems",
			Handler:    _KeeperServiceV1_ListItems_Handler,
		},
		{
			MethodName: "SetItem",
			Handler:    _KeeperServiceV1_SetItem_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "keeper/v1/keeper.proto",
}
