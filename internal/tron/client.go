package tron

import (
    "context"
    "blockchain-wallet/internal/domain"
)

type Client interface {
    CreateWallet(ctx context.Context) (string, error)
    GetBalance(ctx context.Context, address string) (float64, error)
    SendTransaction(ctx context.Context, from, to string, amount float64) (string, error)
    GetTransactionStatus(ctx context.Context, txID string) (*domain.Transaction, error)
}