package grpcpb

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
	WalletService_CreateWallet_FullMethodName    = "/wallet.v1.WalletService/CreateWallet"
	WalletService_GetWallet_FullMethodName       = "/wallet.v1.WalletService/GetWallet"
	WalletService_ListWallets_FullMethodName     = "/wallet.v1.WalletService/ListWallets"
	WalletService_SignMessage_FullMethodName     = "/wallet.v1.WalletService/SignMessage"
	WalletService_SignTransaction_FullMethodName = "/wallet.v1.WalletService/SignTransaction"
	WalletService_GetBalance_FullMethodName      = "/wallet.v1.WalletService/GetBalance"
)

// WalletServiceClient is the client API for WalletService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WalletServiceClient interface {
	CreateWallet(ctx context.Context, in *CreateWalletRequest, opts ...grpc.CallOption) (*WalletResponse, error)
	GetWallet(ctx context.Context, in *GetWalletRequest, opts ...grpc.CallOption) (*WalletResponse, error)
	ListWallets(ctx context.Context, in *ListWalletsRequest, opts ...grpc.CallOption) (*ListWalletsResponse, error)
	SignMessage(ctx context.Context, in *SignMessageRequest, opts ...grpc.CallOption) (*SignMessageResponse, error)
	SignTransaction(ctx context.Context, in *SignTransactionRequest, opts ...grpc.CallOption) (*SignTransactionResponse, error)
	GetBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error)
}

type walletServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWalletServiceClient(cc grpc.ClientConnInterface) WalletServiceClient {
	return &walletServiceClient{cc}
}

func (c *walletServiceClient) CreateWallet(ctx context.Context, in *CreateWalletRequest, opts ...grpc.CallOption) (*WalletResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WalletResponse)
	err := c.cc.Invoke(ctx, WalletService_CreateWallet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) GetWallet(ctx context.Context, in *GetWalletRequest, opts ...grpc.CallOption) (*WalletResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WalletResponse)
	err := c.cc.Invoke(ctx, WalletService_GetWallet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) ListWallets(ctx context.Context, in *ListWalletsRequest, opts ...grpc.CallOption) (*ListWalletsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListWalletsResponse)
	err := c.cc.Invoke(ctx, WalletService_ListWallets_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) SignMessage(ctx context.Context, in *SignMessageRequest, opts ...grpc.CallOption) (*SignMessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SignMessageResponse)
	err := c.cc.Invoke(ctx, WalletService_SignMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) SignTransaction(ctx context.Context, in *SignTransactionRequest, opts ...grpc.CallOption) (*SignTransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SignTransactionResponse)
	err := c.cc.Invoke(ctx, WalletService_SignTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletServiceClient) GetBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetBalanceResponse)
	err := c.cc.Invoke(ctx, WalletService_GetBalance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WalletServiceServer is the server API for WalletService service.
// All implementations must embed UnimplementedWalletServiceServer
// for forward compatibility.
type WalletServiceServer interface {
	CreateWallet(context.Context, *CreateWalletRequest) (*WalletResponse, error)
	GetWallet(context.Context, *GetWalletRequest) (*WalletResponse, error)
	ListWallets(context.Context, *ListWalletsRequest) (*ListWalletsResponse, error)
	SignMessage(context.Context, *SignMessageRequest) (*SignMessageResponse, error)
	SignTransaction(context.Context, *SignTransactionRequest) (*SignTransactionResponse, error)
	GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error)
	mustEmbedUnimplementedWalletServiceServer()
}

// UnimplementedWalletServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedWalletServiceServer struct{}

func (UnimplementedWalletServiceServer) CreateWallet(context.Context, *CreateWalletRequest) (*WalletResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateWallet not implemented")
}
func (UnimplementedWalletServiceServer) GetWallet(context.Context, *GetWalletRequest) (*WalletResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWallet not implemented")
}
func (UnimplementedWalletServiceServer) ListWallets(context.Context, *ListWalletsRequest) (*ListWalletsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListWallets not implemented")
}
func (UnimplementedWalletServiceServer) SignMessage(context.Context, *SignMessageRequest) (*SignMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignMessage not implemented")
}
func (UnimplementedWalletServiceServer) SignTransaction(context.Context, *SignTransactionRequest) (*SignTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignTransaction not implemented")
}
func (UnimplementedWalletServiceServer) GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBalance not implemented")
}
func (UnimplementedWalletServiceServer) mustEmbedUnimplementedWalletServiceServer() {}
func (UnimplementedWalletServiceServer) testEmbeddedByValue()                       {}

// UnsafeWalletServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WalletServiceServer will
// result in compilation errors.
type UnsafeWalletServiceServer interface {
	mustEmbedUnimplementedWalletServiceServer()
}

func RegisterWalletServiceServer(s grpc.ServiceRegistrar, srv WalletServiceServer) {
	// If the following call pancis, it indicates UnimplementedWalletServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&WalletService_ServiceDesc, srv)
}

func _WalletService_CreateWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateWalletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).CreateWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_CreateWallet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).CreateWallet(ctx, req.(*CreateWalletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_GetWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWalletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).GetWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_GetWallet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).GetWallet(ctx, req.(*GetWalletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_ListWallets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListWalletsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).ListWallets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_ListWallets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).ListWallets(ctx, req.(*ListWalletsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_SignMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).SignMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_SignMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).SignMessage(ctx, req.(*SignMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_SignTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).SignTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_SignTransaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).SignTransaction(ctx, req.(*SignTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WalletService_GetBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServiceServer).GetBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WalletService_GetBalance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServiceServer).GetBalance(ctx, req.(*GetBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WalletService_ServiceDesc is the grpc.ServiceDesc for WalletService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WalletService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "wallet.v1.WalletService",
	HandlerType: (*WalletServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateWallet",
			Handler:    _WalletService_CreateWallet_Handler,
		},
		{
			MethodName: "GetWallet",
			Handler:    _WalletService_GetWallet_Handler,
		},
		{
			MethodName: "ListWallets",
			Handler:    _WalletService_ListWallets_Handler,
		},
		{
			MethodName: "SignMessage",
			Handler:    _WalletService_SignMessage_Handler,
		},
		{
			MethodName: "SignTransaction",
			Handler:    _WalletService_SignTransaction_Handler,
		},
		{
			MethodName: "GetBalance",
			Handler:    _WalletService_GetBalance_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/api/grpc/wallet.proto",
}
