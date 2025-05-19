// internal/repository/wallet.go
package repository

import (
    "context"
    "blockchain-wallet/internal/domain"
)

type WalletRepository interface {
    Create(ctx context.Context, wallet *domain.Wallet) error
    FindAll(ctx context.Context, filter domain.WalletFilter) ([]domain.Wallet, domain.Pagination, error)
    FindByAddress(ctx context.Context, address string) (*domain.Wallet, error)
    Update(ctx context.Context, wallet *domain.Wallet) error
}