package domain

import (
	"blockchain-wallet/pkg/blockchain/tron"
	"time"
)

// Transaction представляет транзакцию в блокчейне
// @Description Структура транзакции с детальной информацией о переводе
type Transaction struct {
	Hash          string    `json:"hash" db:"hash" example:"a1b2c3d4e5f6..."`                     // Хеш транзакции в блокчейне
	FromAddress   string    `json:"from_address" db:"from_address" example:"TRX9sGPvkr7i3m1o..."` // Адрес отправителя
	ToAddress     string    `json:"to_address" db:"to_address" example:"TRX8sHQmkc6i4n2p..."`     // Адрес получателя
	Amount        float64   `json:"amount" db:"amount" example:"100.50"`                          // Сумма перевода
	Status        string    `json:"status" db:"status" example:"confirmed"`                       // Статус транзакции (pending, confirmed, failed)
	Confirmations int       `json:"confirmations" db:"confirmations" example:"12"`                // Количество подтверждений
	CreatedAt     time.Time `json:"created_at" db:"created_at" example:"2023-01-01T12:00:00Z"`    // Дата создания
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at" example:"2023-01-01T12:00:00Z"`    // Дата последнего обновления
}

// CreateTransactionRequest запрос для создания транзакции
// @Description Данные, необходимые для создания и отправки транзакции
type CreateTransactionRequest struct {
	FromAddress string         `json:"from_address" binding:"required" example:"TRX9sGPvkr7i3m1o..."` // Адрес кошелька отправителя
	ToAddress   string         `json:"to_address" binding:"required" example:"TRX8sHQmkc6i4n2p..."`   // Адрес кошелька получателя
	Amount      float64        `json:"amount" binding:"required" example:"100.50"`                    // Сумма для перевода
	TokenType   tron.TokenType `json:"token_type" binding:"required,oneof=TRX USDT" example:"TRX"`    // Тип токена: TRX или USDT
}

// TransactionFilter фильтр для поиска транзакций
// @Description Параметры фильтрации при получении списка транзакций
type TransactionFilter struct {
	FromAddress string `json:"from_address" example:"TRX9sGPvkr7i3m1o..."` // Фильтр по адресу отправителя
	ToAddress   string `json:"to_address" example:"TRX8sHQmkc6i4n2p..."`   // Фильтр по адресу получателя
	Status      string `json:"status" example:"confirmed"`                 // Фильтр по статусу транзакции
	Page        int    `json:"page" example:"0"`                           // Номер страницы для пагинации
	Limit       int    `json:"limit" example:"10"`                         // Количество записей на странице
}
