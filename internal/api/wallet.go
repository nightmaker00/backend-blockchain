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

//go:generate mockgen -source=wallet.go -destination=../mocks/wallet_mock.go -package=mocks
type WalletService interface {
	CreateWallet(ctx context.Context, req domain.CreateWalletRequest) (*domain.Wallet, error)
	GetWallets(ctx context.Context, filter domain.WalletFilter) ([]domain.Wallet, domain.Pagination, error)
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, req domain.CreateTransactionRequest) (*domain.Transaction, error)
	GetTransactions(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, domain.Pagination, error)
}

// Объединяющий интерфейс, если требуется
type UseCases interface {
	WalletService
	TransactionService
	// Добавляй другие интерфейсы, если требуется 
}
