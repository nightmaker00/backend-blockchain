package domain

import (
	"time"
)

type Transaction struct {
	TxID          string    `json:"tx_id"`
	FromAddress   string    `json:"from_address"`
	ToAddress     string    `json:"to_address"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	Confirmations int       `json:"confirmations"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateTransactionRequest struct {
	FromAddress string  `json:"from_address"`
	ToAddress   string  `json:"to_address"`
	Amount      float64 `json:"amount"`
}

type TransactionFilter struct {
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Status      string `json:"status"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}

