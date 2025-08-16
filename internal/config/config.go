package config

import (
	"fmt"
	"os"
)

const (
	defaultHTTPPort       = "8080"
	defaultEnv            = "local"
	defaultBaseSepoliaRPC = "https://sepolia.base.org"
	defaultEthSepoliaRPC  = "https://ethereum-sepolia.blockpi.network/v1/rpc/public"
)

type AppConfig struct {
	HTTPPort string
	Env      string

	Networks map[string]NetworkConfig
}

type NetworkConfig struct {
	Name        string
	ChainID     int64
	RPCURL      string
	NativeAsset string
}

func (c *AppConfig) Lookup(key string) (*NetworkConfig, error) {
	network, ok := c.Networks[key]
	if !ok {
		return nil, fmt.Errorf("network %s not configured", key)
	}
	return &network, nil
}

func Load() (*AppConfig, error) {
	cfg := &AppConfig{
		HTTPPort: getEnv("HTTP_PORT", defaultHTTPPort),
		Env:      getEnv("APP_ENV", defaultEnv),
		Networks: map[string]NetworkConfig{
			"base-sepolia": {
				Name:        "Base Sepolia",
				ChainID:     84532,
				RPCURL:      getEnv("BASE_SEPOLIA_RPC_URL", defaultBaseSepoliaRPC),
				NativeAsset: "ETH",
			},
			"eth-sepolia": {
				Name:        "Ethereum Sepolia",
				ChainID:     11155111,
				RPCURL:      getEnv("ETH_SEPOLIA_RPC_URL", defaultEthSepoliaRPC),
				NativeAsset: "ETH",
			},
		},
	}

	if cfg.Networks["base-sepolia"].RPCURL == "" {
		return nil, fmt.Errorf("missing RPC URL for base-sepolia")
	}

	if cfg.Networks["eth-sepolia"].RPCURL == "" {
		return nil, fmt.Errorf("missing RPC URL for eth-sepolia")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
