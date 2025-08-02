package service

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
)

type stubRepo struct {
	mu      sync.RWMutex
	wallets map[string]WalletRecord
}

func newStubRepo() *stubRepo {
	return &stubRepo{wallets: make(map[string]WalletRecord)}
}

func (r *stubRepo) Create(_ context.Context, wallet WalletRecord) (*WalletRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := wallet
	r.wallets[wallet.ID] = copy
	return &copy, nil
}

func (r *stubRepo) GetByID(_ context.Context, id string) (*WalletRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	wallet, ok := r.wallets[id]
	if !ok {
		return nil, ErrNotFound
	}

	copy := wallet
	return &copy, nil
}

func (r *stubRepo) ListByNetwork(_ context.Context, network string) ([]WalletRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]WalletRecord, 0)
	for _, wallet := range r.wallets {
		if network != "" && wallet.Network != network {
			continue
		}
		items = append(items, wallet)
	}
	return items, nil
}

type stubSigner struct {
	newWalletErr    error
	signMessageErr  error
	signTransaction error

	lastNetwork string
	lastPayload []byte
	lastTx      *Transaction
}

func (s *stubSigner) NewWallet(network string) (*WalletRecord, error) {
	s.lastNetwork = network
	if s.newWalletErr != nil {
		return nil, s.newWalletErr
	}
	return &WalletRecord{
		Network:   network,
		Address:   "0x1111111111111111111111111111111111111111",
		PublicKey: "pub-key",
		PrivKey:   "priv-key",
	}, nil
}

func (s *stubSigner) SignMessage(network string, privKeyHex string, payload []byte) (*SignatureOutput, error) {
	s.lastNetwork = network
	s.lastPayload = append([]byte(nil), payload...)
	if s.signMessageErr != nil {
		return nil, s.signMessageErr
	}
	return &SignatureOutput{Signature: "signature", PublicKey: "pub-key"}, nil
}

func (s *stubSigner) SignTransaction(tx *Transaction, privKeyHex string) (string, error) {
	s.lastTx = tx
	if s.signTransaction != nil {
		return "", s.signTransaction
	}
	return "signed-tx", nil
}

func TestCreateWalletRequiresNetwork(t *testing.T) {
	repo := newStubRepo()
	signer := &stubSigner{}
	svc := NewWalletService(repo, signer)

	_, err := svc.CreateWallet(context.Background(), "   ")
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestCreateWalletPersistsRecord(t *testing.T) {
	repo := newStubRepo()
	signer := &stubSigner{}
	svc := NewWalletService(repo, signer)

	wallet, err := svc.CreateWallet(context.Background(), "base-sepolia")
	if err != nil {
		t.Fatalf("CreateWallet returned error: %v", err)
	}
	if wallet.ID == "" {
		t.Fatal("expected wallet ID to be set")
	}
	if wallet.CreatedAt.IsZero() {
		t.Fatal("expected wallet CreatedAt to be set")
	}

	stored, err := repo.GetByID(context.Background(), wallet.ID)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}
	if stored.Network != "base-sepolia" {
		t.Fatalf("expected network base-sepolia, got %s", stored.Network)
	}
}

func TestSignMessageSurfaceNotFound(t *testing.T) {
	repo := newStubRepo()
	signer := &stubSigner{}
	svc := NewWalletService(repo, signer)

	_, err := svc.SignMessage(context.Background(), "missing", []byte("hello"))
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSignTransactionValidatesInput(t *testing.T) {
	repo := newStubRepo()
	signer := &stubSigner{}
	svc := NewWalletService(repo, signer)

	wallet, err := svc.CreateWallet(context.Background(), "eth-sepolia")
	if err != nil {
		t.Fatalf("CreateWallet returned error: %v", err)
	}

	tx := &Transaction{
		ChainID:  11155111,
		To:       "not-an-address",
		Value:    "0x1",
		GasLimit: 21000,
		GasPrice: "0x1",
		Nonce:    1,
	}

	if _, err := svc.SignTransaction(context.Background(), wallet.ID, tx); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestSignTransactionDelegatesToSigner(t *testing.T) {
	repo := newStubRepo()
	signer := &stubSigner{}
	svc := NewWalletService(repo, signer)

	wallet, err := svc.CreateWallet(context.Background(), "eth-sepolia")
	if err != nil {
		t.Fatalf("CreateWallet returned error: %v", err)
	}

	tx := &Transaction{
		ChainID:  11155111,
		To:       "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Value:    "0x1",
		GasLimit: 21000,
		GasPrice: "0x1",
		Nonce:    7,
	}

	signed, err := svc.SignTransaction(context.Background(), wallet.ID, tx)
	if err != nil {
		t.Fatalf("SignTransaction returned error: %v", err)
	}
	if signed != "signed-tx" {
		t.Fatalf("expected signed-tx, got %s", signed)
	}
	if signer.lastTx == nil {
		t.Fatal("expected signer to capture transaction")
	}
	if signer.lastTx.Nonce != 7 {
		t.Fatalf("expected nonce 7, got %d", signer.lastTx.Nonce)
	}
}

func TestListWalletsFiltersByNetwork(t *testing.T) {
	repo := newStubRepo()
	signer := &stubSigner{}
	svc := NewWalletService(repo, signer)

	_, err := svc.CreateWallet(context.Background(), "base-sepolia")
	if err != nil {
		t.Fatalf("CreateWallet returned error: %v", err)
	}
	_, err = svc.CreateWallet(context.Background(), "eth-sepolia")
	if err != nil {
		t.Fatalf("CreateWallet returned error: %v", err)
	}

	wallets, err := svc.ListWallets(context.Background(), "base-sepolia")
	if err != nil {
		t.Fatalf("ListWallets returned error: %v", err)
	}
	if len(wallets) != 1 {
		t.Fatalf("expected 1 wallet, got %d", len(wallets))
	}
	if strings.TrimSpace(wallets[0].Network) != "base-sepolia" {
		t.Fatalf("expected base-sepolia network, got %s", wallets[0].Network)
	}
}

func TestSignMessageUsesSigner(t *testing.T) {
	repo := newStubRepo()
	signer := &stubSigner{}
	svc := NewWalletService(repo, signer)

	wallet, err := svc.CreateWallet(context.Background(), "eth-sepolia")
	if err != nil {
		t.Fatalf("CreateWallet returned error: %v", err)
	}

	msg := []byte("gm")
	output, err := svc.SignMessage(context.Background(), wallet.ID, msg)
	if err != nil {
		t.Fatalf("SignMessage returned error: %v", err)
	}
	if output.Signature == "" {
		t.Fatal("expected signature to be returned")
	}
	if signer.lastNetwork != "eth-sepolia" {
		t.Fatalf("expected signer to receive network eth-sepolia, got %s", signer.lastNetwork)
	}
	if !strings.EqualFold(string(signer.lastPayload), string(msg)) {
		t.Fatal("expected signer to receive payload copy")
	}
}

func TestCreateWalletPropagatesSignerError(t *testing.T) {
	repo := newStubRepo()
	signer := &stubSigner{newWalletErr: errors.New("boom")}
	svc := NewWalletService(repo, signer)

	_, err := svc.CreateWallet(context.Background(), "base-sepolia")
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected propagated signer error, got %v", err)
	}
	if signer.lastNetwork != "base-sepolia" {
		t.Fatalf("expected network to be passed to signer, got %s", signer.lastNetwork)
	}
}

func TestSignTransactionPropagatesSignerError(t *testing.T) {
	repo := newStubRepo()
	signer := &stubSigner{signTransaction: errors.New("sign failed")}
	svc := NewWalletService(repo, signer)

	wallet, err := svc.CreateWallet(context.Background(), "eth-sepolia")
	if err != nil {
		t.Fatalf("CreateWallet returned error: %v", err)
	}

	tx := &Transaction{
		ChainID:  11155111,
		To:       "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		Value:    "0x2",
		GasLimit: 42000,
		GasPrice: "0x5",
		Nonce:    2,
	}

	_, err = svc.SignTransaction(context.Background(), wallet.ID, tx)
	if err == nil || !strings.Contains(err.Error(), "sign failed") {
		t.Fatalf("expected propagated signer error, got %v", err)
	}
	if signer.lastTx == nil {
		t.Fatal("expected signer to receive transaction despite error")
	}
	if signer.lastTx.GasLimit != 42000 {
		t.Fatalf("expected gas limit 42000, got %d", signer.lastTx.GasLimit)
	}
}

func TestGetWalletSetsTimestamps(t *testing.T) {
	repo := newStubRepo()
	signer := &stubSigner{}
	svc := NewWalletService(repo, signer)

	wallet, err := svc.CreateWallet(context.Background(), "eth-sepolia")
	if err != nil {
		t.Fatalf("CreateWallet returned error: %v", err)
	}

	fetched, err := svc.GetWallet(context.Background(), wallet.ID)
	if err != nil {
		t.Fatalf("GetWallet returned error: %v", err)
	}
	if fetched.CreatedAt.Before(time.Now().Add(-time.Minute)) {
		t.Fatalf("expected recent CreatedAt timestamp, got %v", fetched.CreatedAt)
	}
}
