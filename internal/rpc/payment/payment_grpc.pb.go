// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package payment

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

// PaymentsClient is the client API for Payments service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PaymentsClient interface {
	MakePayment(ctx context.Context, in *PaymentRequest, opts ...grpc.CallOption) (*PaymentResponse, error)
	PaymentStatus(ctx context.Context, in *PaymentStatusRequest, opts ...grpc.CallOption) (*PaymentStatusResponse, error)
}

type paymentsClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentsClient(cc grpc.ClientConnInterface) PaymentsClient {
	return &paymentsClient{cc}
}

func (c *paymentsClient) MakePayment(ctx context.Context, in *PaymentRequest, opts ...grpc.CallOption) (*PaymentResponse, error) {
	out := new(PaymentResponse)
	err := c.cc.Invoke(ctx, "/rpc.Payments/MakePayment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentsClient) PaymentStatus(ctx context.Context, in *PaymentStatusRequest, opts ...grpc.CallOption) (*PaymentStatusResponse, error) {
	out := new(PaymentStatusResponse)
	err := c.cc.Invoke(ctx, "/rpc.Payments/PaymentStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentsServer is the server API for Payments service.
// All implementations must embed UnimplementedPaymentsServer
// for forward compatibility
type PaymentsServer interface {
	MakePayment(context.Context, *PaymentRequest) (*PaymentResponse, error)
	PaymentStatus(context.Context, *PaymentStatusRequest) (*PaymentStatusResponse, error)
	mustEmbedUnimplementedPaymentsServer()
}

// UnimplementedPaymentsServer must be embedded to have forward compatible implementations.
type UnimplementedPaymentsServer struct {
}

func (UnimplementedPaymentsServer) MakePayment(context.Context, *PaymentRequest) (*PaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakePayment not implemented")
}
func (UnimplementedPaymentsServer) PaymentStatus(context.Context, *PaymentStatusRequest) (*PaymentStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PaymentStatus not implemented")
}
func (UnimplementedPaymentsServer) mustEmbedUnimplementedPaymentsServer() {}

// UnsafePaymentsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaymentsServer will
// result in compilation errors.
type UnsafePaymentsServer interface {
	mustEmbedUnimplementedPaymentsServer()
}

func RegisterPaymentsServer(s grpc.ServiceRegistrar, srv PaymentsServer) {
	s.RegisterService(&Payments_ServiceDesc, srv)
}

func _Payments_MakePayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServer).MakePayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Payments/MakePayment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServer).MakePayment(ctx, req.(*PaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payments_PaymentStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PaymentStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentsServer).PaymentStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.Payments/PaymentStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentsServer).PaymentStatus(ctx, req.(*PaymentStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Payments_ServiceDesc is the grpc.ServiceDesc for Payments service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Payments_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.Payments",
	HandlerType: (*PaymentsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MakePayment",
			Handler:    _Payments_MakePayment_Handler,
		},
		{
			MethodName: "PaymentStatus",
			Handler:    _Payments_PaymentStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "payment/payment.proto",
}