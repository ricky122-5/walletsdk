package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/rickyreddygari/walletsdk/internal/testutil"
)

func TestWalletLifecycle(t *testing.T) {
	server, cleanup := testutil.NewTestServer(t)
	defer cleanup()

	client := server.Client()

	payload := map[string]string{"network": "base-sepolia"}
	body, _ := json.Marshal(payload)

	resp := testutil.MustDo(t, client, mustRequest(t, http.MethodPost, server.URL+"/v1/wallets", bytes.NewReader(body)))
	defer resp.Body.Close()

	testutil.AssertStatus(t, resp, http.StatusCreated)

	var wallet struct {
		ID        string `json:"id"`
		Network   string `json:"network"`
		Address   string `json:"address"`
		PublicKey string `json:"publicKey"`
	}
	testutil.DecodeJSON(t, resp, &wallet)

	if wallet.Network != "base-sepolia" {
		t.Fatalf("expected network base-sepolia, got %s", wallet.Network)
	}

	resp = testutil.MustDo(t, client, mustRequest(t, http.MethodGet, fmt.Sprintf("%s/v1/wallets/%s", server.URL, wallet.ID), nil))
	defer resp.Body.Close()
	testutil.AssertStatus(t, resp, http.StatusOK)

	var fetched map[string]interface{}
	testutil.DecodeJSON(t, resp, &fetched)

	resp = testutil.MustDo(t, client, mustRequest(t, http.MethodGet, fmt.Sprintf("%s/v1/wallets?network=%s", server.URL, "base-sepolia"), nil))
	defer resp.Body.Close()
	testutil.AssertStatus(t, resp, http.StatusOK)

	var wallets []map[string]interface{}
	testutil.DecodeJSON(t, resp, &wallets)
	if len(wallets) != 1 {
		t.Fatalf("expected 1 wallet, got %d", len(wallets))
	}

	messagePayload := map[string]string{"message": "gm"}
	msgBody, _ := json.Marshal(messagePayload)
	resp = testutil.MustDo(t, client, mustRequest(t, http.MethodPost, fmt.Sprintf("%s/v1/wallets/%s/sign-message", server.URL, wallet.ID), bytes.NewReader(msgBody)))
	defer resp.Body.Close()
	testutil.AssertStatus(t, resp, http.StatusOK)

	var signature struct {
		Signature string `json:"signature"`
		PublicKey string `json:"publicKey"`
	}
	testutil.DecodeJSON(t, resp, &signature)
	if signature.Signature == "" {
		t.Fatalf("expected signature to be non-empty")
	}
}

func mustRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}
