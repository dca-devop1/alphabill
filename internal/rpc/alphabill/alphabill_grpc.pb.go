// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// AlphaBillServiceClient is the client API for AlphaBillService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AlphaBillServiceClient interface {
	GetBlocks(ctx context.Context, in *GetBlocksRequest, opts ...grpc.CallOption) (AlphaBillService_GetBlocksClient, error)
	ProcessTransaction(ctx context.Context, in *transaction.Transaction, opts ...grpc.CallOption) (*transaction.TransactionResponse, error)
}

type alphaBillServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAlphaBillServiceClient(cc grpc.ClientConnInterface) AlphaBillServiceClient {
	return &alphaBillServiceClient{cc}
}

func (c *alphaBillServiceClient) GetBlocks(ctx context.Context, in *GetBlocksRequest, opts ...grpc.CallOption) (AlphaBillService_GetBlocksClient, error) {
	stream, err := c.cc.NewStream(ctx, &AlphaBillService_ServiceDesc.Streams[0], "/abrpc.AlphaBillService/GetBlocks", opts...)
	if err != nil {
		return nil, err
	}
	x := &alphaBillServiceGetBlocksClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type AlphaBillService_GetBlocksClient interface {
	Recv() (*GetBlocksResponse, error)
	grpc.ClientStream
}

type alphaBillServiceGetBlocksClient struct {
	grpc.ClientStream
}

func (x *alphaBillServiceGetBlocksClient) Recv() (*GetBlocksResponse, error) {
	m := new(GetBlocksResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *alphaBillServiceClient) ProcessTransaction(ctx context.Context, in *transaction.Transaction, opts ...grpc.CallOption) (*transaction.TransactionResponse, error) {
	out := new(transaction.TransactionResponse)
	err := c.cc.Invoke(ctx, "/abrpc.AlphaBillService/ProcessTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AlphaBillServiceServer is the server API for AlphaBillService service.
// All implementations must embed UnimplementedAlphaBillServiceServer
// for forward compatibility
type AlphaBillServiceServer interface {
	GetBlocks(*GetBlocksRequest, AlphaBillService_GetBlocksServer) error
	ProcessTransaction(context.Context, *transaction.Transaction) (*transaction.TransactionResponse, error)
	mustEmbedUnimplementedAlphaBillServiceServer()
}

// UnimplementedAlphaBillServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAlphaBillServiceServer struct {
}

func (UnimplementedAlphaBillServiceServer) GetBlocks(*GetBlocksRequest, AlphaBillService_GetBlocksServer) error {
	return status.Errorf(codes.Unimplemented, "method GetBlocks not implemented")
}
func (UnimplementedAlphaBillServiceServer) ProcessTransaction(context.Context, *transaction.Transaction) (*transaction.TransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessTransaction not implemented")
}
func (UnimplementedAlphaBillServiceServer) mustEmbedUnimplementedAlphaBillServiceServer() {}

// UnsafeAlphaBillServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AlphaBillServiceServer will
// result in compilation errors.
type UnsafeAlphaBillServiceServer interface {
	mustEmbedUnimplementedAlphaBillServiceServer()
}

func RegisterAlphaBillServiceServer(s grpc.ServiceRegistrar, srv AlphaBillServiceServer) {
	s.RegisterService(&AlphaBillService_ServiceDesc, srv)
}

func _AlphaBillService_GetBlocks_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetBlocksRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AlphaBillServiceServer).GetBlocks(m, &alphaBillServiceGetBlocksServer{stream})
}

type AlphaBillService_GetBlocksServer interface {
	Send(*GetBlocksResponse) error
	grpc.ServerStream
}

type alphaBillServiceGetBlocksServer struct {
	grpc.ServerStream
}

func (x *alphaBillServiceGetBlocksServer) Send(m *GetBlocksResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _AlphaBillService_ProcessTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(transaction.Transaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlphaBillServiceServer).ProcessTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/abrpc.AlphaBillService/ProcessTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlphaBillServiceServer).ProcessTransaction(ctx, req.(*transaction.Transaction))
	}
	return interceptor(ctx, in, info, handler)
}

// AlphaBillService_ServiceDesc is the grpc.ServiceDesc for AlphaBillService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AlphaBillService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "abrpc.AlphaBillService",
	HandlerType: (*AlphaBillServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProcessTransaction",
			Handler:    _AlphaBillService_ProcessTransaction_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetBlocks",
			Handler:       _AlphaBillService_GetBlocks_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "alphabill.proto",
}
