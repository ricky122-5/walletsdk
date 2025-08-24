package grpc_test

import (
	"context"
	"testing"

	grpcpb "github.com/rickyreddygari/walletsdk/internal/api/grpcpb"
	"github.com/rickyreddygari/walletsdk/internal/testutil"
)

func TestGRPCWalletLifecycle(t *testing.T) {
	_, conn, cleanup := testutil.NewTestServer(t)
	defer cleanup()

	client := grpcpb.NewWalletServiceClient(conn)

	wallet := testutil.MustInvoke(t, context.Background(), func(ctx context.Context) (*grpcpb.WalletResponse, error) {
		return client.CreateWallet(ctx, &grpcpb.CreateWalletRequest{Network: "base-sepolia"})
	})

	fetched := testutil.MustInvoke(t, context.Background(), func(ctx context.Context) (*grpcpb.WalletResponse, error) {
		return client.GetWallet(ctx, &grpcpb.GetWalletRequest{WalletId: wallet.Id})
	})
	if fetched.Id != wallet.Id {
		t.Fatalf("expected wallet id %s, got %s", wallet.Id, fetched.Id)
	}

	list := testutil.MustInvoke(t, context.Background(), func(ctx context.Context) (*grpcpb.ListWalletsResponse, error) {
		return client.ListWallets(ctx, &grpcpb.ListWalletsRequest{Network: "base-sepolia"})
	})
	if len(list.Wallets) == 0 {
		t.Fatalf("expected wallets in list")
	}

	sig := testutil.MustInvoke(t, context.Background(), func(ctx context.Context) (*grpcpb.SignMessageResponse, error) {
		return client.SignMessage(ctx, &grpcpb.SignMessageRequest{WalletId: wallet.Id, Payload: []byte("gm")})
	})
	if sig.Signature == "" {
		t.Fatalf("expected signature")
	}

	balance := testutil.MustInvoke(t, context.Background(), func(ctx context.Context) (*grpcpb.GetBalanceResponse, error) {
		return client.GetBalance(ctx, &grpcpb.GetBalanceRequest{WalletId: wallet.Id})
	})
	if balance.Balance == nil {
		t.Fatalf("expected balance")
	}
}

