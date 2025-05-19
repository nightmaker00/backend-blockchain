// internal/domain/models.go
package domain

import "time"

type WalletType string

const (
    WalletTypeRegular WalletType = "regular"
    WalletTypeBank    WalletType = "bank"
)

type Wallet struct {
    Address    string     `json:"address"`
    WalletType WalletType `json:"wallet_type"`
    Name       string     `json:"name"`
    Status     string     `json:"status"`
    Balance    float64    `json:"balance"`
    CreatedAt  time.Time  `json:"created_at"`
}

type CreateWalletRequest struct {
    WalletType WalletType `json:"wallet_type" validate:"required,oneof=regular bank"`
    Name       string     `json:"name"`
}

type WalletFilter struct {
    WalletType string
    Status     string
    Page       int
    Limit      int
}

type Pagination struct {
    Total int `json:"total"`
    Page  int `json:"page"`
    Limit int `json:"limit"`
}

type WalletsResponse struct {
    Wallets    []Wallet   `json:"wallets"`
    Pagination Pagination `json:"pagination"`
}