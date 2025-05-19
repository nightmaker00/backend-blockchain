package domain

import (
	"context"
)

// WalletService определяет контракт для работы с кошельками
type WalletService interface {
	CreateWallet(ctx context.Context, req CreateWalletRequest) (*Wallet, error)
	GetWallets(ctx context.Context, filter WalletFilter) ([]Wallet, Pagination, error)
}

// ... existing code ...
