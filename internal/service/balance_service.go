package service

import (
	"context"
	"fmt"
)

// BalanceRepository access persisted wallets for balance checks.
type BalanceRepository interface {
	GetByID(ctx context.Context, id string) (*WalletRecord, error)
}

// BalanceFetcher fetches balances from networks.
type BalanceFetcher interface {
	FetchBalance(ctx context.Context, rpcURL string, address string) (string, error)
}

// NetworkRegistry fetches network metadata.
type NetworkRegistry interface {
	Lookup(network string) (*Network, error)
}

// Network describes an accessible blockchain network.
type Network struct {
	Name        string
	ChainID     int64
	RPCURL      string
	NativeAsset string
}

type balanceService struct {
	repo     BalanceRepository
	fetcher  BalanceFetcher
	registry NetworkRegistry
}

// BalanceService exposes balance operations.
type BalanceService interface {
	GetBalance(ctx context.Context, walletID string) (*Balance, error)
}

// NewBalanceService constructs the balance service.
func NewBalanceService(repo BalanceRepository, fetcher BalanceFetcher, registry NetworkRegistry) BalanceService {
	return &balanceService{
		repo:     repo,
		fetcher:  fetcher,
		registry: registry,
	}
}

// GetBalance fetches the balance for a wallet using its network RPC.
func (s *balanceService) GetBalance(ctx context.Context, walletID string) (*Balance, error) {
	record, err := s.repo.GetByID(ctx, walletID)
	if err != nil {
		if err == ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get wallet: %w", err)
	}

	network, err := s.registry.Lookup(record.Network)
	if err != nil {
		return nil, fmt.Errorf("lookup network: %w", err)
	}

	amount, err := s.fetcher.FetchBalance(ctx, network.RPCURL, record.Address)
	if err != nil {
		return nil, fmt.Errorf("fetch balance: %w", err)
	}

	return &Balance{Asset: network.NativeAsset, Amount: amount}, nil
}
