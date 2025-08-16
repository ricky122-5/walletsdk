package app

import (
	"fmt"

	httprouter "github.com/rickyreddygari/walletsdk/internal/api/http"
	"github.com/rickyreddygari/walletsdk/internal/blockchain/ethereum"
	"github.com/rickyreddygari/walletsdk/internal/config"
	"github.com/rickyreddygari/walletsdk/internal/service"
	"github.com/rickyreddygari/walletsdk/internal/storage/memory"
)

type Container struct {
	Config         *config.AppConfig
	WalletService  service.WalletService
	BalanceService service.BalanceService
	HTTPServer     *httprouter.Server
}

func NewContainer() (*Container, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	repo := memory.NewWalletRepository()
	signer := ethereum.NewSigner()
	fetcher := ethereum.NewBalanceFetcher()
	registry := service.NewConfigRegistry(cfg)

	walletService := service.NewWalletService(repo, signer)
	balanceService := service.NewBalanceService(repo, fetcher, registry)

	httpServer := httprouter.NewServer()
	routes := httprouter.NewRouteBuilder(walletService, balanceService)
	routes.Register(httpServer.Router())

	return &Container{
		Config:         cfg,
		WalletService:  walletService,
		BalanceService: balanceService,
		HTTPServer:     httpServer,
	}, nil
}
