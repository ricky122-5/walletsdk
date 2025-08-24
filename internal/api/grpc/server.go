package grpc

import (
	"context"

	"github.com/rickyreddygari/walletsdk/internal/api/grpcpb"
	"github.com/rickyreddygari/walletsdk/internal/service"
)

type Server struct {
	grpcpb.UnimplementedWalletServiceServer
	wallets  service.WalletService
	balances service.BalanceService
}

func NewServer(wallets service.WalletService, balances service.BalanceService) *Server {
	return &Server{wallets: wallets, balances: balances}
}

func (s *Server) CreateWallet(ctx context.Context, req *grpcpb.CreateWalletRequest) (*grpcpb.WalletResponse, error) {
	wallet, err := s.wallets.CreateWallet(ctx, req.GetNetwork())
	if err != nil {
		return nil, err
	}
	return toProtoWallet(wallet), nil
}

func (s *Server) GetWallet(ctx context.Context, req *grpcpb.GetWalletRequest) (*grpcpb.WalletResponse, error) {
	wallet, err := s.wallets.GetWallet(ctx, req.GetWalletId())
	if err != nil {
		return nil, err
	}
	return toProtoWallet(wallet), nil
}

func (s *Server) ListWallets(ctx context.Context, req *grpcpb.ListWalletsRequest) (*grpcpb.ListWalletsResponse, error) {
	wallets, err := s.wallets.ListWallets(ctx, req.GetNetwork())
	if err != nil {
		return nil, err
	}

	resp := &grpcpb.ListWalletsResponse{Wallets: make([]*grpcpb.WalletResponse, 0, len(wallets))}
	for _, wallet := range wallets {
		w := wallet
		resp.Wallets = append(resp.Wallets, toProtoWallet(&w))
	}
	return resp, nil
}

func (s *Server) SignMessage(ctx context.Context, req *grpcpb.SignMessageRequest) (*grpcpb.SignMessageResponse, error) {
	sig, err := s.wallets.SignMessage(ctx, req.GetWalletId(), req.GetPayload())
	if err != nil {
		return nil, err
	}
	return &grpcpb.SignMessageResponse{Signature: sig.Signature, PublicKey: sig.PublicKey}, nil
}

func (s *Server) SignTransaction(ctx context.Context, req *grpcpb.SignTransactionRequest) (*grpcpb.SignTransactionResponse, error) {
	tx := fromProtoTransaction(req.GetTransaction())
	signed, err := s.wallets.SignTransaction(ctx, req.GetWalletId(), tx)
	if err != nil {
		return nil, err
	}
	return &grpcpb.SignTransactionResponse{SignedTransaction: signed}, nil
}

func (s *Server) GetBalance(ctx context.Context, req *grpcpb.GetBalanceRequest) (*grpcpb.GetBalanceResponse, error) {
	balance, err := s.balances.GetBalance(ctx, req.GetWalletId())
	if err != nil {
		return nil, err
	}
	return &grpcpb.GetBalanceResponse{Balance: &grpcpb.Balance{Asset: balance.Asset, Amount: balance.Amount}}, nil
}

func toProtoWallet(wallet *service.Wallet) *grpcpb.WalletResponse {
	var createdAt int64
	if !wallet.CreatedAt.IsZero() {
		createdAt = wallet.CreatedAt.Unix()
	}
	return &grpcpb.WalletResponse{
		Id:            wallet.ID,
		Network:       wallet.Network,
		Address:       wallet.Address,
		PublicKey:     wallet.PublicKey,
		CreatedAtUnix: createdAt,
	}
}

func fromProtoTransaction(tx *grpcpb.Transaction) *service.Transaction {
	if tx == nil {
		return nil
	}
	return &service.Transaction{
		ChainID:  tx.GetChainId(),
		From:     tx.GetFrom(),
		To:       tx.GetTo(),
		Value:    tx.GetValue(),
		Data:     tx.GetData(),
		GasLimit: tx.GetGasLimit(),
		GasPrice: tx.GetGasPrice(),
		Nonce:    tx.GetNonce(),
	}
}
