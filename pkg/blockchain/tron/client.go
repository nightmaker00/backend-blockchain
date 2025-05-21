package tron

import (
	"context"
	"net/http"
)

type Client interface {
	CreateWallet(ctx context.Context) (string, error)
	GetBalance(ctx context.Context, address string) (float64, error)
	GetTransactionStatus(ctx context.Context, txID string) (*Transaction, error)
	SendTransaction(ctx context.Context, fromAddress, toAddress string, amount float64) (string, error)
}

type TronCLient struct {
	httpClient *http.Client
}

func NewTronClient(httpClient *http.Client) *TronCLient {
	return &TronCLient{
		httpClient: httpClient,
	}
}
