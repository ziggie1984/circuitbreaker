// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package circuitbreakerrpc

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

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	GetInfo(ctx context.Context, in *GetInfoRequest, opts ...grpc.CallOption) (*GetInfoResponse, error)
	UpdateLimit(ctx context.Context, in *UpdateLimitRequest, opts ...grpc.CallOption) (*UpdateLimitResponse, error)
	// Clear specific limit and use default.
	ClearLimit(ctx context.Context, in *ClearLimitRequest, opts ...grpc.CallOption) (*ClearLimitResponse, error)
	UpdateDefaultLimit(ctx context.Context, in *UpdateDefaultLimitRequest, opts ...grpc.CallOption) (*UpdateDefaultLimitResponse, error)
	ListLimits(ctx context.Context, in *ListLimitsRequest, opts ...grpc.CallOption) (*ListLimitsResponse, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) GetInfo(ctx context.Context, in *GetInfoRequest, opts ...grpc.CallOption) (*GetInfoResponse, error) {
	out := new(GetInfoResponse)
	err := c.cc.Invoke(ctx, "/circuitbreaker.Service/GetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) UpdateLimit(ctx context.Context, in *UpdateLimitRequest, opts ...grpc.CallOption) (*UpdateLimitResponse, error) {
	out := new(UpdateLimitResponse)
	err := c.cc.Invoke(ctx, "/circuitbreaker.Service/UpdateLimit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ClearLimit(ctx context.Context, in *ClearLimitRequest, opts ...grpc.CallOption) (*ClearLimitResponse, error) {
	out := new(ClearLimitResponse)
	err := c.cc.Invoke(ctx, "/circuitbreaker.Service/ClearLimit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) UpdateDefaultLimit(ctx context.Context, in *UpdateDefaultLimitRequest, opts ...grpc.CallOption) (*UpdateDefaultLimitResponse, error) {
	out := new(UpdateDefaultLimitResponse)
	err := c.cc.Invoke(ctx, "/circuitbreaker.Service/UpdateDefaultLimit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) ListLimits(ctx context.Context, in *ListLimitsRequest, opts ...grpc.CallOption) (*ListLimitsResponse, error) {
	out := new(ListLimitsResponse)
	err := c.cc.Invoke(ctx, "/circuitbreaker.Service/ListLimits", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	GetInfo(context.Context, *GetInfoRequest) (*GetInfoResponse, error)
	UpdateLimit(context.Context, *UpdateLimitRequest) (*UpdateLimitResponse, error)
	// Clear specific limit and use default.
	ClearLimit(context.Context, *ClearLimitRequest) (*ClearLimitResponse, error)
	UpdateDefaultLimit(context.Context, *UpdateDefaultLimitRequest) (*UpdateDefaultLimitResponse, error)
	ListLimits(context.Context, *ListLimitsRequest) (*ListLimitsResponse, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) GetInfo(context.Context, *GetInfoRequest) (*GetInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInfo not implemented")
}
func (UnimplementedServiceServer) UpdateLimit(context.Context, *UpdateLimitRequest) (*UpdateLimitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLimit not implemented")
}
func (UnimplementedServiceServer) ClearLimit(context.Context, *ClearLimitRequest) (*ClearLimitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearLimit not implemented")
}
func (UnimplementedServiceServer) UpdateDefaultLimit(context.Context, *UpdateDefaultLimitRequest) (*UpdateDefaultLimitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDefaultLimit not implemented")
}
func (UnimplementedServiceServer) ListLimits(context.Context, *ListLimitsRequest) (*ListLimitsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLimits not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_GetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/circuitbreaker.Service/GetInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetInfo(ctx, req.(*GetInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_UpdateLimit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateLimitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).UpdateLimit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/circuitbreaker.Service/UpdateLimit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).UpdateLimit(ctx, req.(*UpdateLimitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ClearLimit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearLimitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ClearLimit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/circuitbreaker.Service/ClearLimit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ClearLimit(ctx, req.(*ClearLimitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_UpdateDefaultLimit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDefaultLimitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).UpdateDefaultLimit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/circuitbreaker.Service/UpdateDefaultLimit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).UpdateDefaultLimit(ctx, req.(*UpdateDefaultLimitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_ListLimits_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLimitsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ListLimits(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/circuitbreaker.Service/ListLimits",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ListLimits(ctx, req.(*ListLimitsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "circuitbreaker.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetInfo",
			Handler:    _Service_GetInfo_Handler,
		},
		{
			MethodName: "UpdateLimit",
			Handler:    _Service_UpdateLimit_Handler,
		},
		{
			MethodName: "ClearLimit",
			Handler:    _Service_ClearLimit_Handler,
		},
		{
			MethodName: "UpdateDefaultLimit",
			Handler:    _Service_UpdateDefaultLimit_Handler,
		},
		{
			MethodName: "ListLimits",
			Handler:    _Service_ListLimits_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "circuitbreaker.proto",
}
