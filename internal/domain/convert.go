package domain

import (
	"blockchain-wallet/pkg/blockchain/tron"
	"time"
)

func ToDomainTransaction(tx *tron.Transaction) *Transaction {
	if len(tx.RawData.Contract) == 0 {
		return nil
	}
	value := tx.RawData.Contract[0].Parameter.Value

	return &Transaction{
		Hash:          tx.TxID,
		FromAddress:   value.OwnerAddress,
		ToAddress:     value.ToAddress,
		Amount:        float64(value.Amount) / 1_000_000,
		Status:        "pending",
		Confirmations: 0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
