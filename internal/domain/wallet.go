package domain

import (
	"context"
	"time"
)

type Wallet struct {
	Address    string    `json:"address"`
	WalletType string    `json:"wallet_type"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateWalletRequest struct {
	WalletType string `json:"wallet_type"`
	Name       string `json:"name"`
}

type WalletFilter struct {
	WalletType string `json:"wallet_type"`
	Status     string `json:"status"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type WalletService interface {
	CreateWallet(ctx context.Context, req CreateWalletRequest) (*Wallet, error)
	GetWallets(ctx context.Context, filter WalletFilter) ([]Wallet, Pagination, error)
}

// ... existing code ...
