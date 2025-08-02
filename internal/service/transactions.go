package service

// Transaction represents a simplified transaction request.
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
