// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.1.0
// - protoc             v3.15.7
// source: proto/service.proto

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

// McocServiceClient is the client API for McocService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type McocServiceClient interface {
	AddChamp(ctx context.Context, in *AddChampRequest, opts ...grpc.CallOption) (*AddChampResponse, error)
	LockChamp(ctx context.Context, in *LockChampRequest, opts ...grpc.CallOption) (*LockChampResponse, error)
	ListChamps(ctx context.Context, in *ListChampsRequest, opts ...grpc.CallOption) (*ListChampsResponse, error)
	GetWarDefense(ctx context.Context, in *GetWarDefenseRequest, opts ...grpc.CallOption) (*GetWarDefenseResponse, error)
	UpdateChamp(ctx context.Context, in *AddChampRequest, opts ...grpc.CallOption) (*AddChampResponse, error)
}

type mcocServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMcocServiceClient(cc grpc.ClientConnInterface) McocServiceClient {
	return &mcocServiceClient{cc}
}

func (c *mcocServiceClient) AddChamp(ctx context.Context, in *AddChampRequest, opts ...grpc.CallOption) (*AddChampResponse, error) {
	out := new(AddChampResponse)
	err := c.cc.Invoke(ctx, "/proto.McocService/AddChamp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mcocServiceClient) LockChamp(ctx context.Context, in *LockChampRequest, opts ...grpc.CallOption) (*LockChampResponse, error) {
	out := new(LockChampResponse)
	err := c.cc.Invoke(ctx, "/proto.McocService/LockChamp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mcocServiceClient) ListChamps(ctx context.Context, in *ListChampsRequest, opts ...grpc.CallOption) (*ListChampsResponse, error) {
	out := new(ListChampsResponse)
	err := c.cc.Invoke(ctx, "/proto.McocService/ListChamps", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mcocServiceClient) GetWarDefense(ctx context.Context, in *GetWarDefenseRequest, opts ...grpc.CallOption) (*GetWarDefenseResponse, error) {
	out := new(GetWarDefenseResponse)
	err := c.cc.Invoke(ctx, "/proto.McocService/GetWarDefense", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mcocServiceClient) UpdateChamp(ctx context.Context, in *AddChampRequest, opts ...grpc.CallOption) (*AddChampResponse, error) {
	out := new(AddChampResponse)
	err := c.cc.Invoke(ctx, "/proto.McocService/UpdateChamp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// McocServiceServer is the server API for McocService service.
// All implementations must embed UnimplementedMcocServiceServer
// for forward compatibility
type McocServiceServer interface {
	AddChamp(context.Context, *AddChampRequest) (*AddChampResponse, error)
	LockChamp(context.Context, *LockChampRequest) (*LockChampResponse, error)
	ListChamps(context.Context, *ListChampsRequest) (*ListChampsResponse, error)
	GetWarDefense(context.Context, *GetWarDefenseRequest) (*GetWarDefenseResponse, error)
	UpdateChamp(context.Context, *AddChampRequest) (*AddChampResponse, error)
	mustEmbedUnimplementedMcocServiceServer()
}

// UnimplementedMcocServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMcocServiceServer struct {
}

func (UnimplementedMcocServiceServer) AddChamp(context.Context, *AddChampRequest) (*AddChampResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddChamp not implemented")
}
func (UnimplementedMcocServiceServer) LockChamp(context.Context, *LockChampRequest) (*LockChampResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LockChamp not implemented")
}
func (UnimplementedMcocServiceServer) ListChamps(context.Context, *ListChampsRequest) (*ListChampsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListChamps not implemented")
}
func (UnimplementedMcocServiceServer) GetWarDefense(context.Context, *GetWarDefenseRequest) (*GetWarDefenseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWarDefense not implemented")
}
func (UnimplementedMcocServiceServer) UpdateChamp(context.Context, *AddChampRequest) (*AddChampResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateChamp not implemented")
}
func (UnimplementedMcocServiceServer) mustEmbedUnimplementedMcocServiceServer() {}

// UnsafeMcocServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to McocServiceServer will
// result in compilation errors.
type UnsafeMcocServiceServer interface {
	mustEmbedUnimplementedMcocServiceServer()
}

func RegisterMcocServiceServer(s grpc.ServiceRegistrar, srv McocServiceServer) {
	s.RegisterService(&McocService_ServiceDesc, srv)
}

func _McocService_AddChamp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddChampRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McocServiceServer).AddChamp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.McocService/AddChamp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McocServiceServer).AddChamp(ctx, req.(*AddChampRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _McocService_LockChamp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LockChampRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McocServiceServer).LockChamp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.McocService/LockChamp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McocServiceServer).LockChamp(ctx, req.(*LockChampRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _McocService_ListChamps_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListChampsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McocServiceServer).ListChamps(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.McocService/ListChamps",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McocServiceServer).ListChamps(ctx, req.(*ListChampsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _McocService_GetWarDefense_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWarDefenseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McocServiceServer).GetWarDefense(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.McocService/GetWarDefense",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McocServiceServer).GetWarDefense(ctx, req.(*GetWarDefenseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _McocService_UpdateChamp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddChampRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McocServiceServer).UpdateChamp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.McocService/UpdateChamp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McocServiceServer).UpdateChamp(ctx, req.(*AddChampRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// McocService_ServiceDesc is the grpc.ServiceDesc for McocService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var McocService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.McocService",
	HandlerType: (*McocServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddChamp",
			Handler:    _McocService_AddChamp_Handler,
		},
		{
			MethodName: "LockChamp",
			Handler:    _McocService_LockChamp_Handler,
		},
		{
			MethodName: "ListChamps",
			Handler:    _McocService_ListChamps_Handler,
		},
		{
			MethodName: "GetWarDefense",
			Handler:    _McocService_GetWarDefense_Handler,
		},
		{
			MethodName: "UpdateChamp",
			Handler:    _McocService_UpdateChamp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/service.proto",
}
