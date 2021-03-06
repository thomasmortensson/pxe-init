// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// PxeInitClient is the client API for PxeInit service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PxeInitClient interface {
	// List images
	ListImages(ctx context.Context, in *ListImagesRequest, opts ...grpc.CallOption) (*ListImagesResponse, error)
	// Register image to machine with MAC
	RegisterImageMachine(ctx context.Context, in *RegisterImageMachineRequest, opts ...grpc.CallOption) (*RegisterImageMachineResponse, error)
}

type pxeInitClient struct {
	cc grpc.ClientConnInterface
}

func NewPxeInitClient(cc grpc.ClientConnInterface) PxeInitClient {
	return &pxeInitClient{cc}
}

func (c *pxeInitClient) ListImages(ctx context.Context, in *ListImagesRequest, opts ...grpc.CallOption) (*ListImagesResponse, error) {
	out := new(ListImagesResponse)
	err := c.cc.Invoke(ctx, "/v1.PxeInit/ListImages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pxeInitClient) RegisterImageMachine(ctx context.Context, in *RegisterImageMachineRequest, opts ...grpc.CallOption) (*RegisterImageMachineResponse, error) {
	out := new(RegisterImageMachineResponse)
	err := c.cc.Invoke(ctx, "/v1.PxeInit/RegisterImageMachine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PxeInitServer is the server API for PxeInit service.
// All implementations must embed UnimplementedPxeInitServer
// for forward compatibility
type PxeInitServer interface {
	// List images
	ListImages(context.Context, *ListImagesRequest) (*ListImagesResponse, error)
	// Register image to machine with MAC
	RegisterImageMachine(context.Context, *RegisterImageMachineRequest) (*RegisterImageMachineResponse, error)
	mustEmbedUnimplementedPxeInitServer()
}

// UnimplementedPxeInitServer must be embedded to have forward compatible implementations.
type UnimplementedPxeInitServer struct {
}

func (UnimplementedPxeInitServer) ListImages(context.Context, *ListImagesRequest) (*ListImagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListImages not implemented")
}
func (UnimplementedPxeInitServer) RegisterImageMachine(context.Context, *RegisterImageMachineRequest) (*RegisterImageMachineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterImageMachine not implemented")
}
func (UnimplementedPxeInitServer) mustEmbedUnimplementedPxeInitServer() {}

// UnsafePxeInitServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PxeInitServer will
// result in compilation errors.
type UnsafePxeInitServer interface {
	mustEmbedUnimplementedPxeInitServer()
}

func RegisterPxeInitServer(s grpc.ServiceRegistrar, srv PxeInitServer) {
	s.RegisterService(&PxeInit_ServiceDesc, srv)
}

func _PxeInit_ListImages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListImagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PxeInitServer).ListImages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.PxeInit/ListImages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PxeInitServer).ListImages(ctx, req.(*ListImagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PxeInit_RegisterImageMachine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterImageMachineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PxeInitServer).RegisterImageMachine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.PxeInit/RegisterImageMachine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PxeInitServer).RegisterImageMachine(ctx, req.(*RegisterImageMachineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PxeInit_ServiceDesc is the grpc.ServiceDesc for PxeInit service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PxeInit_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.PxeInit",
	HandlerType: (*PxeInitServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListImages",
			Handler:    _PxeInit_ListImages_Handler,
		},
		{
			MethodName: "RegisterImageMachine",
			Handler:    _PxeInit_RegisterImageMachine_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pxe_init.proto",
}
