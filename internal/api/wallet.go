package api

import (
	"blockchain-wallet/internal/domain"
	"context"
)

// WalletService определяет контракт для работы с кошельками
type WalletService interface {
	CreateWallet(ctx context.Context, req domain.CreateWalletRequest) (*domain.Wallet, error)
	GetWallets(ctx context.Context, filter domain.WalletFilter) ([]domain.Wallet, domain.Pagination, error)
}
