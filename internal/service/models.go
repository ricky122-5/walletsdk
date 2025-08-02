package service

import "time"

type Wallet struct {
	ID        string
	Network   string
	Address   string
	PublicKey string
	CreatedAt time.Time
}

type Balance struct {
	Asset  string
	Amount string
}

type SignatureOutput struct {
	Signature string
	PublicKey string
}
