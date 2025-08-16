package service

import (
	"context"
	"fmt"
)

type BalanceRepository interface {
	GetByID(ctx context.Context, id string) (*WalletRecord, error)
}

type BalanceFetcher interface {
	FetchBalance(ctx context.Context, rpcURL string, address string) (string, error)
}

type NetworkRegistry interface {
	Lookup(network string) (*Network, error)
}

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

type BalanceService interface {
	GetBalance(ctx context.Context, walletID string) (*Balance, error)
}

func NewBalanceService(repo BalanceRepository, fetcher BalanceFetcher, registry NetworkRegistry) BalanceService {
	return &balanceService{
		repo:     repo,
		fetcher:  fetcher,
		registry: registry,
	}
}

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
