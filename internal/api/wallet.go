package api

import (
	"blockchain-wallet/internal/domain"
	"context"
)

// SOLID
// S - Single Responsibility Principle
// O - Open/Closed Principle
// L - Liskov Substitution Principle
// I - Interface Segregation Principle
// D - Dependency Inversion Principle

type WalletService interface {
	CreateWallet(ctx context.Context, req domain.CreateWalletRequest) (*domain.Wallet, error)
	GetBalance(ctx context.Context, address string) (float64, error)
	SendTransaction(ctx context.Context, req domain.CreateTransactionRequest) (*domain.Transaction, error)
	GetTransactionStatus(ctx context.Context, txID string) (string, error)
	GetTransactions(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, domain.Pagination, error)
	GetWallets(ctx context.Context, filter domain.WalletFilter) ([]domain.Wallet, domain.Pagination, error)
}
