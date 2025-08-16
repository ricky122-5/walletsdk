package ethereum

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BalanceFetcher struct {
	clientFactory func(rpcURL string) (*ethclient.Client, error)
}

func NewBalanceFetcher() *BalanceFetcher {
	return &BalanceFetcher{
		clientFactory: ethclient.Dial,
	}
}

func (f *BalanceFetcher) WithClientFactory(factory func(rpcURL string) (*ethclient.Client, error)) {
	f.clientFactory = factory
}

func (f *BalanceFetcher) FetchBalance(ctx context.Context, rpcURL string, address string) (string, error) {
	client, err := f.clientFactory(rpcURL)
	if err != nil {
		return "", fmt.Errorf("dial rpc: %w", err)
	}
	defer client.Close()

	addr := common.HexToAddress(address)
	balance, err := client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return "", fmt.Errorf("fetch balance: %w", err)
	}

	return formatWei(balance), nil
}

func formatWei(v *big.Int) string {
	if v == nil {
		return "0"
	}
	return v.String()
}
