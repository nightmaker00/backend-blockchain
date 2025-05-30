package service

import (
	tronlib "blockchain-wallet/pkg/blockchain/tron"
	"context"
)

type TronClient interface {
	GetBalance(ctx context.Context, address string) (float64, error)
	SendTransaction(ctx context.Context, fromAddress, toAddress string, amount float64, privateKey string) (*tronlib.Transaction, error)
	GetTransactionStatus(ctx context.Context, txID string) (string, error)
}
