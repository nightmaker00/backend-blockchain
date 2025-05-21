package domain

import (
	"time"
)

type WalletType string

const (
	WalletTypeRegular WalletType = "regular"
	WalletTypeBank    WalletType = "bank"
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

type WalletsResponse struct {
	Wallets    []Wallet   `json:"wallets"`
	Pagination Pagination `json:"pagination"`
}
