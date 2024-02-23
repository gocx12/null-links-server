// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.6.1
// source: rpc_service/content_security/content_security.proto

package content_security

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
	ContentSecurityService_WebsetCheckQuery_FullMethodName   = "/content_security.ContentSecurityService/WebsetCheckQuery"
	ContentSecurityService_WebsetCheckConfirm_FullMethodName = "/content_security.ContentSecurityService/WebsetCheckConfirm"
)

// ContentSecurityServiceClient is the client API for ContentSecurityService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ContentSecurityServiceClient interface {
	WebsetCheckQuery(ctx context.Context, in *WebsetCheckQueryReq, opts ...grpc.CallOption) (*WebsetCheckQueryResp, error)
	WebsetCheckConfirm(ctx context.Context, in *WebsetCheckConfirmReq, opts ...grpc.CallOption) (*WebsetCheckConfirmResp, error)
}

type contentSecurityServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewContentSecurityServiceClient(cc grpc.ClientConnInterface) ContentSecurityServiceClient {
	return &contentSecurityServiceClient{cc}
}

func (c *contentSecurityServiceClient) WebsetCheckQuery(ctx context.Context, in *WebsetCheckQueryReq, opts ...grpc.CallOption) (*WebsetCheckQueryResp, error) {
	out := new(WebsetCheckQueryResp)
	err := c.cc.Invoke(ctx, ContentSecurityService_WebsetCheckQuery_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentSecurityServiceClient) WebsetCheckConfirm(ctx context.Context, in *WebsetCheckConfirmReq, opts ...grpc.CallOption) (*WebsetCheckConfirmResp, error) {
	out := new(WebsetCheckConfirmResp)
	err := c.cc.Invoke(ctx, ContentSecurityService_WebsetCheckConfirm_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ContentSecurityServiceServer is the server API for ContentSecurityService service.
// All implementations must embed UnimplementedContentSecurityServiceServer
// for forward compatibility
type ContentSecurityServiceServer interface {
	WebsetCheckQuery(context.Context, *WebsetCheckQueryReq) (*WebsetCheckQueryResp, error)
	WebsetCheckConfirm(context.Context, *WebsetCheckConfirmReq) (*WebsetCheckConfirmResp, error)
	mustEmbedUnimplementedContentSecurityServiceServer()
}

// UnimplementedContentSecurityServiceServer must be embedded to have forward compatible implementations.
type UnimplementedContentSecurityServiceServer struct {
}

func (UnimplementedContentSecurityServiceServer) WebsetCheckQuery(context.Context, *WebsetCheckQueryReq) (*WebsetCheckQueryResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WebsetCheckQuery not implemented")
}
func (UnimplementedContentSecurityServiceServer) WebsetCheckConfirm(context.Context, *WebsetCheckConfirmReq) (*WebsetCheckConfirmResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WebsetCheckConfirm not implemented")
}
func (UnimplementedContentSecurityServiceServer) mustEmbedUnimplementedContentSecurityServiceServer() {
}

// UnsafeContentSecurityServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ContentSecurityServiceServer will
// result in compilation errors.
type UnsafeContentSecurityServiceServer interface {
	mustEmbedUnimplementedContentSecurityServiceServer()
}

func RegisterContentSecurityServiceServer(s grpc.ServiceRegistrar, srv ContentSecurityServiceServer) {
	s.RegisterService(&ContentSecurityService_ServiceDesc, srv)
}

func _ContentSecurityService_WebsetCheckQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WebsetCheckQueryReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentSecurityServiceServer).WebsetCheckQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentSecurityService_WebsetCheckQuery_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentSecurityServiceServer).WebsetCheckQuery(ctx, req.(*WebsetCheckQueryReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ContentSecurityService_WebsetCheckConfirm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WebsetCheckConfirmReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentSecurityServiceServer).WebsetCheckConfirm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ContentSecurityService_WebsetCheckConfirm_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentSecurityServiceServer).WebsetCheckConfirm(ctx, req.(*WebsetCheckConfirmReq))
	}
	return interceptor(ctx, in, info, handler)
}

// ContentSecurityService_ServiceDesc is the grpc.ServiceDesc for ContentSecurityService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ContentSecurityService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "content_security.ContentSecurityService",
	HandlerType: (*ContentSecurityServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WebsetCheckQuery",
			Handler:    _ContentSecurityService_WebsetCheckQuery_Handler,
		},
		{
			MethodName: "WebsetCheckConfirm",
			Handler:    _ContentSecurityService_WebsetCheckConfirm_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc_service/content_security/content_security.proto",
}