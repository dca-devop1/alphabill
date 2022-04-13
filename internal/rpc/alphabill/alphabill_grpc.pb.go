// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.11.2
// source: alphabill.proto

package alphabill

import (
	context "context"
	transaction "gitdc.ee.guardtime.com/alphabill/alphabill/internal/rpc/transaction"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AlphabillServiceClient is the client API for AlphabillService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AlphabillServiceClient interface {
	ProcessTransaction(ctx context.Context, in *transaction.Transaction, opts ...grpc.CallOption) (*transaction.TransactionResponse, error)
	GetBlock(ctx context.Context, in *GetBlockRequest, opts ...grpc.CallOption) (*GetBlockResponse, error)
	GetMaxBlockNo(ctx context.Context, in *GetMaxBlockNoRequest, opts ...grpc.CallOption) (*GetMaxBlockNoResponse, error)
}

type alphabillServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAlphabillServiceClient(cc grpc.ClientConnInterface) AlphabillServiceClient {
	return &alphabillServiceClient{cc}
}

func (c *alphabillServiceClient) ProcessTransaction(ctx context.Context, in *transaction.Transaction, opts ...grpc.CallOption) (*transaction.TransactionResponse, error) {
	out := new(transaction.TransactionResponse)
	err := c.cc.Invoke(ctx, "/abrpc.AlphabillService/ProcessTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alphabillServiceClient) GetBlock(ctx context.Context, in *GetBlockRequest, opts ...grpc.CallOption) (*GetBlockResponse, error) {
	out := new(GetBlockResponse)
	err := c.cc.Invoke(ctx, "/abrpc.AlphabillService/GetBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alphabillServiceClient) GetMaxBlockNo(ctx context.Context, in *GetMaxBlockNoRequest, opts ...grpc.CallOption) (*GetMaxBlockNoResponse, error) {
	out := new(GetMaxBlockNoResponse)
	err := c.cc.Invoke(ctx, "/abrpc.AlphabillService/GetMaxBlockNo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AlphabillServiceServer is the server API for AlphabillService service.
// All implementations must embed UnimplementedAlphabillServiceServer
// for forward compatibility
type AlphabillServiceServer interface {
	ProcessTransaction(context.Context, *transaction.Transaction) (*transaction.TransactionResponse, error)
	GetBlock(context.Context, *GetBlockRequest) (*GetBlockResponse, error)
	GetMaxBlockNo(context.Context, *GetMaxBlockNoRequest) (*GetMaxBlockNoResponse, error)
	mustEmbedUnimplementedAlphabillServiceServer()
}

// UnimplementedAlphabillServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAlphabillServiceServer struct {
}

func (UnimplementedAlphabillServiceServer) ProcessTransaction(context.Context, *transaction.Transaction) (*transaction.TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessTransaction not implemented")
}
func (UnimplementedAlphabillServiceServer) GetBlock(context.Context, *GetBlockRequest) (*GetBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlock not implemented")
}
func (UnimplementedAlphabillServiceServer) GetMaxBlockNo(context.Context, *GetMaxBlockNoRequest) (*GetMaxBlockNoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMaxBlockNo not implemented")
}
func (UnimplementedAlphabillServiceServer) mustEmbedUnimplementedAlphabillServiceServer() {}

// UnsafeAlphabillServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AlphabillServiceServer will
// result in compilation errors.
type UnsafeAlphabillServiceServer interface {
	mustEmbedUnimplementedAlphabillServiceServer()
}

func RegisterAlphabillServiceServer(s grpc.ServiceRegistrar, srv AlphabillServiceServer) {
	s.RegisterService(&AlphabillService_ServiceDesc, srv)
}

func _AlphabillService_ProcessTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(transaction.Transaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlphabillServiceServer).ProcessTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/abrpc.AlphabillService/ProcessTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlphabillServiceServer).ProcessTransaction(ctx, req.(*transaction.Transaction))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlphabillService_GetBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlphabillServiceServer).GetBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/abrpc.AlphabillService/GetBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlphabillServiceServer).GetBlock(ctx, req.(*GetBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlphabillService_GetMaxBlockNo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMaxBlockNoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlphabillServiceServer).GetMaxBlockNo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/abrpc.AlphabillService/GetMaxBlockNo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlphabillServiceServer).GetMaxBlockNo(ctx, req.(*GetMaxBlockNoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AlphabillService_ServiceDesc is the grpc.ServiceDesc for AlphabillService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AlphabillService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "abrpc.AlphabillService",
	HandlerType: (*AlphabillServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProcessTransaction",
			Handler:    _AlphabillService_ProcessTransaction_Handler,
		},
		{
			MethodName: "GetBlock",
			Handler:    _AlphabillService_GetBlock_Handler,
		},
		{
			MethodName: "GetMaxBlockNo",
			Handler:    _AlphabillService_GetMaxBlockNo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "alphabill.proto",
}
