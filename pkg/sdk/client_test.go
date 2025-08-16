package sdk_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rickyreddygari/walletsdk/internal/app"
	"github.com/rickyreddygari/walletsdk/pkg/sdk"
)

func TestClientWalletLifecycle(t *testing.T) {
	container, err := app.NewContainer()
	if err != nil {
		t.Fatalf("bootstrap container: %v", err)
	}

	server := httptest.NewServer(container.HTTPServer)
	defer server.Close()

	client, err := sdk.NewClient(server.URL)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	wallet, err := client.CreateWallet(sdk.CreateWalletRequest{Network: "base-sepolia"})
	if err != nil {
		t.Fatalf("CreateWallet failed: %v", err)
	}

	fetched, err := client.GetWallet(wallet.ID)
	if err != nil {
		t.Fatalf("GetWallet failed: %v", err)
	}
	if fetched.ID != wallet.ID {
		t.Fatalf("expected wallet ID %s, got %s", wallet.ID, fetched.ID)
	}

	wallets, err := client.ListWallets("base-sepolia")
	if err != nil {
		t.Fatalf("ListWallets failed: %v", err)
	}
	if len(wallets) == 0 {
		t.Fatalf("expected at least one wallet in list")
	}

	signature, err := client.SignMessage(wallet.ID, sdk.SignMessageRequest{Message: "gm"})
	if err != nil {
		t.Fatalf("SignMessage failed: %v", err)
	}
	if signature.Signature == "" {
		t.Fatalf("expected signature to be set")
	}

	balance, err := client.GetBalance(wallet.ID)
	if err != nil {
		t.Fatalf("GetBalance failed: %v", err)
	}
	if balance.Asset == "" {
		t.Fatalf("expected balance asset to be set")
	}
}

type mockRoundTrip struct {
	status int
	body   map[string]string
}

func (m *mockRoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: m.status,
		Body:       io.NopCloser(bytes.NewBuffer(nil)),
		Header:     make(http.Header),
	}
	if m.body != nil {
		encoded, _ := json.Marshal(m.body)
		resp.Body = io.NopCloser(bytes.NewReader(encoded))
	}
	return resp, nil
}

func TestClientHandlesAPIError(t *testing.T) {
	custom := &http.Client{Transport: &mockRoundTrip{status: http.StatusBadRequest, body: map[string]string{"error": "bad request"}}}
	client, err := sdk.NewClient("https://example.com", sdk.WithHTTPClient(custom))
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	_, err = client.CreateWallet(sdk.CreateWalletRequest{Network: "base-sepolia"})
	if err == nil {
		t.Fatalf("expected error from API")
	}
}
