package service

import (
	"fmt"

	"github.com/rickyreddygari/walletsdk/internal/config"
)

type ConfigRegistry struct {
	cfg *config.AppConfig
}

func NewConfigRegistry(cfg *config.AppConfig) *ConfigRegistry {
	return &ConfigRegistry{cfg: cfg}
}

func (r *ConfigRegistry) Lookup(network string) (*Network, error) {
	if network == "" {
		return nil, fmt.Errorf("%w: network required", ErrValidation)
	}

	cfg, err := r.cfg.Lookup(network)
	if err != nil {
		return nil, err
	}

	return &Network{
		Name:        cfg.Name,
		ChainID:     cfg.ChainID,
		RPCURL:      cfg.RPCURL,
		NativeAsset: cfg.NativeAsset,
	}, nil
}
