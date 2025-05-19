package service

import (
    "context"
    "blockchain-wallet/internal/domain"
)

type TransactionService interface {
    // Пример методов, которые могут понадобиться
    GetTransactionStatus(ctx context.Context, txID string) (*domain.Transaction, error)
    // Добавьте другие методы по необходимости
}