package domain

import (
	"time"
)

type WalletKind string

const (
	WalletKindRegular WalletKind = "regular"
	WalletKindBank    WalletKind = "bank"
)

type Wallet struct {
	PublicKey  string     `json:"public_key" db:"public_key"`
	PrivateKey string     `json:"private_key" db:"private_key"`
	Address    string     `json:"address" db:"address"`
	SeedPhrase string     `json:"seed_phrase" db:"seed_phrase"`
	Kind       WalletKind `json:"kind" db:"kind"`
	IsActive   bool       `json:"is_active" db:"is_active"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
	Username   string     `json:"username" db:"username"`
}

type CreateWalletRequest struct {
	Kind     string `json:"kind"`
	Username string `json:"username"`
}

type WalletFilter struct {
	Kind     string `json:"kind"`
	IsActive bool   `json:"is_active"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type WalletsResponse struct {
	Wallets    []Wallet   `json:"wallets"`
	Pagination Pagination `json:"pagination"`
}

