package ethereum

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"

	"github.com/rickyreddygari/walletsdk/internal/service"
)

type Signer struct{}

func NewSigner() *Signer {
	return &Signer{}
}

func (s *Signer) NewWallet(network string) (*service.WalletRecord, error) {
	key, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	address := crypto.PubkeyToAddress(key.PublicKey)
	privBytes := crypto.FromECDSA(key)

	return &service.WalletRecord{
		Network:   network,
		Address:   address.Hex(),
		PublicKey: hex.EncodeToString(crypto.FromECDSAPub(&key.PublicKey)),
		PrivKey:   hex.EncodeToString(privBytes),
	}, nil
}

func (s *Signer) SignMessage(_ string, privKeyHex string, payload []byte) (*service.SignatureOutput, error) {
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		return nil, fmt.Errorf("decode private key: %w", err)
	}

	digest := hashMessage(payload)

	sig, err := crypto.Sign(digest, privKey)
	if err != nil {
		return nil, fmt.Errorf("sign digest: %w", err)
	}

	if len(sig) != 65 {
		return nil, fmt.Errorf("unexpected signature length: %d", len(sig))
	}

	sig[64] += 27

	return &service.SignatureOutput{
		Signature: "0x" + hex.EncodeToString(sig),
		PublicKey: "0x" + hex.EncodeToString(crypto.FromECDSAPub(&privKey.PublicKey)),
	}, nil
}

func hashMessage(payload []byte) []byte {
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(payload))
	digest := sha3.NewLegacyKeccak256()
	digest.Write([]byte(prefix))
	digest.Write(payload)
	return digest.Sum(nil)
}

func RecoverAddress(message []byte, signature []byte) (common.Address, error) {
	if len(signature) != 65 {
		return common.Address{}, fmt.Errorf("invalid signature length: %d", len(signature))
	}

	sig := make([]byte, len(signature))
	copy(sig, signature)

	if sig[64] >= 27 {
		sig[64] -= 27
	}

	digest := hashMessage(message)

	pubKey, err := crypto.SigToPub(digest, sig)
	if err != nil {
		return common.Address{}, fmt.Errorf("recover public key: %w", err)
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}

func (s *Signer) SignTransaction(tx *service.Transaction, privKeyHex string) (string, error) {
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		return "", fmt.Errorf("decode private key: %w", err)
	}

	to := common.HexToAddress(tx.To)
	value, ok := new(big.Int).SetString(stripHex(tx.Value), 16)
	if !ok {
		return "", fmt.Errorf("parse value")
	}

	gasPrice, ok := new(big.Int).SetString(stripHex(tx.GasPrice), 16)
	if !ok {
		return "", fmt.Errorf("parse gas price")
	}

	data := common.FromHex(tx.Data)

	baseTx := types.NewTransaction(
		tx.Nonce,
		to,
		value,
		tx.GasLimit,
		gasPrice,
		data,
	)

	signer := types.NewEIP155Signer(big.NewInt(tx.ChainID))
	signedTx, err := types.SignTx(baseTx, signer, privKey)
	if err != nil {
		return "", fmt.Errorf("sign tx: %w", err)
	}

	bytes, err := signedTx.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("marshal signed tx: %w", err)
	}

	return "0x" + hex.EncodeToString(bytes), nil
}

func stripHex(input string) string {
	if len(input) >= 2 && input[:2] == "0x" {
		return input[2:]
	}
	return input
}
