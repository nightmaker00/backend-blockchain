package service

import (
	"context"

	"blockchain-wallet/internal/domain"
	"blockchain-wallet/pkg/blockchain/tron"
)

type TronClient interface {
	CreateWallet(ctx context.Context) (string, error)
	GetBalance(ctx context.Context, address string) (float64, error)
	SendTransaction(ctx context.Context, from, to string, amount float64) (string, error)
	GetTransactionStatus(ctx context.Context, txID string) (*tron.Transaction, error)
}

func ToDomainTransaction(tx *tron.Transaction) *domain.Transaction {
	return &domain.Transaction{
		TxID: tx.TxID,
		FromAddress: tx.FromAddress,
		ToAddress: tx.ToAddress,
		Amount: tx.Amount,
		Status: tx.Status,
		Confirmations: tx.Confirmations,
		CreatedAt: tx.CreatedAt,
	}
}
