package domain

import (
	"blockchain-wallet/pkg/blockchain/tron"
	"time"
)

type Transaction struct {
	Hash          string    `json:"hash" db:"hash"`
	FromAddress   string    `json:"from_address" db:"from_address"`
	ToAddress     string    `json:"to_address" db:"to_address"`
	Amount        float64   `json:"amount" db:"amount"`
	Status        string    `json:"status" db:"status"`
	Confirmations int       `json:"confirmations" db:"confirmations"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTransactionRequest struct {
	FromAddress string         `json:"from_address" binding:"required"`
	ToAddress   string         `json:"to_address" binding:"required"`
	Amount      float64        `json:"amount" binding:"required"`
	TokenType   tron.TokenType `json:"token_type" binding:"required,oneof=TRX USDT"`
}

type TransactionFilter struct {
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Status      string `json:"status"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}
