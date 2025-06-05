package service

import (
	"blockchain-wallet/internal/domain"
	"context"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *domain.Wallet) error
	FindAll(ctx context.Context, filter domain.WalletFilter) ([]domain.Wallet, domain.Pagination, error)
	FindByAddress(ctx context.Context, address string) (*domain.Wallet, error)
	Update(ctx context.Context, wallet *domain.Wallet) error
	GetTransactions(ctx context.Context, address string, pagination *domain.Pagination) ([]domain.Transaction, domain.Pagination, error)
	GetTransactionStatus(ctx context.Context, txID string) (string, error)
	SaveTransaction(ctx context.Context, trx domain.Transaction) error
}
