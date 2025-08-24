package testutil

import (
	"context"
	"testing"

	"google.golang.org/grpc"

	grpcpb "github.com/rickyreddygari/walletsdk/internal/api/grpcpb"
)

func NewGRPCClient(t *testing.T, conn *grpc.ClientConn) grpcpb.WalletServiceClient {
	t.Helper()
	return grpcpb.NewWalletServiceClient(conn)
}

func MustInvoke[T any](t *testing.T, ctx context.Context, fn func(context.Context) (T, error)) T {
	t.Helper()
	resp, err := fn(ctx)
	if err != nil {
		t.Fatalf("grpc call failed: %v", err)
	}
	return resp
}
