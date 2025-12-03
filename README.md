# Wallet SDK Service

Wallet-as-a-Service style backend that generates wallets, signs payloads, and exposes a Go SDK for builders. This is a Demo project, wouldn't recommend using in prod before changing the in-memory storage system + priv key management. Happy building!

## Quick Start

```bash
go run ./cmd/server
```

```bash
curl -X POST http://localhost:8080/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{"network":"base-sepolia"}'
```

### Docker

```bash
docker build -t walletsdk:latest .
docker run --rm -p 8080:8080 walletsdk:latest
```

### AWS Lambda

```bash
sam local start-api
```

### Go SDK Usage

```go
client, _ := sdk.NewClient("http://localhost:8080")
wallet, _ := client.CreateWallet(sdk.CreateWalletRequest{Network: "base-sepolia"})
signature, _ := client.SignMessage(wallet.ID, sdk.SignMessageRequest{Message: "gm"})
```

### Tests

```bash
go test ./...
```
