package service

import (
	tronlib "blockchain-wallet/pkg/blockchain/tron"
	"context"
)

type TronClient interface {
	GetBalance(ctx context.Context, address string) (*tronlib.WalletBalance, error)
	SendToken(ctx context.Context, fromAddress, toAddress string, amount float64, privateKey string, tokenType tronlib.TokenType) (interface{}, error)
}
