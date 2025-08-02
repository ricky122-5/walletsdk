package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents the wallet SDK client.
type Client struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

// NewClient constructs the SDK client with default settings.
func NewClient(baseURL string, opts ...Option) (*Client, error) {
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	client := &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
	for _, opt := range opts {
		opt(client)
	}
	return client, nil
}

// Option configures Client instances.
type Option func(*Client)

// WithHTTPClient overrides the default HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithAPIKey sets an API key header for requests.
func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = key
	}
}

// CreateWalletRequest represents the payload for wallet creation.
type CreateWalletRequest struct {
	Network string `json:"network"`
}

// WalletResponse represents the wallet data returned from the service.
type WalletResponse struct {
	ID        string    `json:"id"`
	Network   string    `json:"network"`
	Address   string    `json:"address"`
	PublicKey string    `json:"publicKey"`
	CreatedAt time.Time `json:"createdAt"`
}

// CreateWallet creates a wallet on the service.
func (c *Client) CreateWallet(req CreateWalletRequest) (*WalletResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/wallets", c.baseURL), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiErr map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("api error (%d)", resp.StatusCode)
		}
		return nil, fmt.Errorf("api error (%d): %s", resp.StatusCode, apiErr["error"])
	}

	var wallet WalletResponse
	if err := json.NewDecoder(resp.Body).Decode(&wallet); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &wallet, nil
}

// WalletListResponse wraps the array response from the service.
type WalletListResponse []WalletResponse

// SignatureResponse represents the signing output.
type SignatureResponse struct {
	Signature string `json:"signature"`
	PublicKey string `json:"publicKey"`
}

// BalanceResponse captures the balance payload.
type BalanceResponse struct {
	Asset  string `json:"asset"`
	Amount string `json:"amount"`
}

// SignMessageRequest represents a signing request.
type SignMessageRequest struct {
	Message string `json:"message"`
}

// Transaction mirrors the service transaction payload.
type Transaction struct {
	ChainID  int64  `json:"chainId"`
	From     string `json:"from,omitempty"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data,omitempty"`
	GasLimit uint64 `json:"gasLimit"`
	GasPrice string `json:"gasPrice"`
	Nonce    uint64 `json:"nonce"`
}

func (c *Client) GetWallet(id string) (*WalletResponse, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("wallet id is required")
	}

	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("%s/v1/wallets/%s", c.baseURL, id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wallet WalletResponse
	if err := json.NewDecoder(resp.Body).Decode(&wallet); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &wallet, nil
}

func (c *Client) ListWallets(network string) (WalletListResponse, error) {
	endpoint := fmt.Sprintf("%s/v1/wallets", c.baseURL)
	if network != "" {
		endpoint = fmt.Sprintf("%s?v=1&network=%s", endpoint, url.QueryEscape(network))
	}

	resp, err := c.doRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wallets WalletListResponse
	if err := json.NewDecoder(resp.Body).Decode(&wallets); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return wallets, nil
}

func (c *Client) SignMessage(walletID string, req SignMessageRequest) (*SignatureResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(http.MethodPost, fmt.Sprintf("%s/v1/wallets/%s/sign-message", c.baseURL, walletID), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var signature SignatureResponse
	if err := json.NewDecoder(resp.Body).Decode(&signature); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &signature, nil
}

func (c *Client) SignTransaction(walletID string, tx *Transaction) (string, error) {
	payload, err := json.Marshal(tx)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.doRequest(http.MethodPost, fmt.Sprintf("%s/v1/wallets/%s/sign-transaction", c.baseURL, walletID), bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body struct {
		SignedTransaction string `json:"signedTransaction"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	return body.SignedTransaction, nil
}

func (c *Client) GetBalance(walletID string) (*BalanceResponse, error) {
	resp, err := c.doRequest(http.MethodGet, fmt.Sprintf("%s/v1/wallets/%s/balance", c.baseURL, walletID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var balance BalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &balance, nil
}

func (c *Client) doRequest(method string, endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		var apiErr map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return nil, fmt.Errorf("api error (%d)", resp.StatusCode)
		}
		return nil, fmt.Errorf("api error (%d): %s", resp.StatusCode, apiErr["error"])
	}

	return resp, nil
}
