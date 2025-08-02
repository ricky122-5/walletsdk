package memory

import (
	"context"
	"errors"
	"sync"

	servicepkg "github.com/rickyreddygari/walletsdk/internal/service"
)

type WalletRepository struct {
	mu      sync.RWMutex
	wallets map[string]servicepkg.WalletRecord
}

func NewWalletRepository() *WalletRepository {
	return &WalletRepository{
		wallets: make(map[string]servicepkg.WalletRecord),
	}
}

func (r *WalletRepository) Create(_ context.Context, wallet servicepkg.WalletRecord) (*servicepkg.WalletRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.wallets[wallet.ID]; exists {
		return nil, errors.New("wallet already exists")
	}

	r.wallets[wallet.ID] = wallet
	copy := wallet
	return &copy, nil
}

func (r *WalletRepository) GetByID(_ context.Context, id string) (*servicepkg.WalletRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	wallet, ok := r.wallets[id]
	if !ok {
		return nil, servicepkg.ErrNotFound
	}
	return &wallet, nil
}

func (r *WalletRepository) ListByNetwork(_ context.Context, network string) ([]servicepkg.WalletRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]servicepkg.WalletRecord, 0)
	for _, wallet := range r.wallets {
		if network != "" && wallet.Network != network {
			continue
		}
		result = append(result, wallet)
	}
	return result, nil
}
