package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet WalletRecord) (*WalletRecord, error)
	GetByID(ctx context.Context, id string) (*WalletRecord, error)
	ListByNetwork(ctx context.Context, network string) ([]WalletRecord, error)
}

type Signer interface {
	NewWallet(network string) (*WalletRecord, error)
	SignMessage(network string, privKeyHex string, payload []byte) (*SignatureOutput, error)
	SignTransaction(tx *Transaction, privKeyHex string) (string, error)
}

type WalletRecord struct {
	ID        string
	Network   string
	Address   string
	PublicKey string
	PrivKey   string
	CreatedAt time.Time
}

type walletService struct {
	repo   WalletRepository
	signer Signer
}

func NewWalletService(repo WalletRepository, signer Signer) WalletService {
	return &walletService{repo: repo, signer: signer}
}

type WalletService interface {
	CreateWallet(ctx context.Context, network string) (*Wallet, error)
	GetWallet(ctx context.Context, id string) (*Wallet, error)
	ListWallets(ctx context.Context, network string) ([]Wallet, error)
	SignMessage(ctx context.Context, walletID string, payload []byte) (*SignatureOutput, error)
	SignTransaction(ctx context.Context, walletID string, tx *Transaction) (string, error)
}

func (s *walletService) CreateWallet(ctx context.Context, network string) (*Wallet, error) {
	network = strings.TrimSpace(network)
	if network == "" {
		return nil, fmt.Errorf("%w: network is required", ErrValidation)
	}

	record, err := s.signer.NewWallet(network)
	if err != nil {
		return nil, fmt.Errorf("generate wallet: %w", err)
	}

	record.ID = uuid.NewString()
	record.CreatedAt = time.Now().UTC()

	stored, err := s.repo.Create(ctx, *record)
	if err != nil {
		return nil, fmt.Errorf("store wallet: %w", err)
	}

	return &Wallet{
		ID:        stored.ID,
		Network:   stored.Network,
		Address:   stored.Address,
		PublicKey: stored.PublicKey,
		CreatedAt: stored.CreatedAt,
	}, nil
}

func (s *walletService) GetWallet(ctx context.Context, id string) (*Wallet, error) {
	record, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get wallet: %w", err)
	}

	return &Wallet{
		ID:        record.ID,
		Network:   record.Network,
		Address:   record.Address,
		PublicKey: record.PublicKey,
		CreatedAt: record.CreatedAt,
	}, nil
}

func (s *walletService) ListWallets(ctx context.Context, network string) ([]Wallet, error) {
	records, err := s.repo.ListByNetwork(ctx, network)
	if err != nil {
		return nil, fmt.Errorf("list wallets: %w", err)
	}

	wallets := make([]Wallet, 0, len(records))
	for _, record := range records {
		wallets = append(wallets, Wallet{
			ID:        record.ID,
			Network:   record.Network,
			Address:   record.Address,
			PublicKey: record.PublicKey,
			CreatedAt: record.CreatedAt,
		})
	}

	return wallets, nil
}

func (s *walletService) SignMessage(ctx context.Context, walletID string, payload []byte) (*SignatureOutput, error) {
	record, err := s.repo.GetByID(ctx, walletID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get wallet: %w", err)
	}

	signature, err := s.signer.SignMessage(record.Network, record.PrivKey, payload)
	if err != nil {
		return nil, fmt.Errorf("sign payload: %w", err)
	}

	return signature, nil
}

func (s *walletService) SignTransaction(ctx context.Context, walletID string, tx *Transaction) (string, error) {
	record, err := s.repo.GetByID(ctx, walletID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("get wallet: %w", err)
	}

	if err := ValidateTransaction(tx); err != nil {
		return "", err
	}

	signed, err := s.signer.SignTransaction(tx, record.PrivKey)
	if err != nil {
		if errors.Is(err, ErrNotImplemented) {
			return "", ErrNotImplemented
		}
		return "", fmt.Errorf("sign transaction: %w", err)
	}

	return signed, nil
}
