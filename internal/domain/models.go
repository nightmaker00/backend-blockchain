package domain

import (
	"time"
)

type WalletKind string

const (
	WalletKindRegular WalletKind = "regular"
	WalletKindBank    WalletKind = "bank"
)

// Wallet представляет блокчейн кошелек в системе
// @Description Структура блокчейн кошелька с ключами и метаданными
type Wallet struct {
	PublicKey  string     `json:"public_key" db:"public_key" example:"0x1234567890abcdef..."`   // Публичный ключ кошелька
	PrivateKey string     `json:"private_key" db:"private_key" example:"0xabcdef1234567890..."` // Приватный ключ кошелька (конфиденциально)
	Address    string     `json:"address" db:"address" example:"TRX9sGPvkr7i3m1o..."`           // Адрес кошелька в сети TRON
	SeedPhrase string     `json:"seed_phrase" db:"seed_phrase" example:"word1 word2 word3..."`  // Мнемоническая фраза для восстановления
	Kind       WalletKind `json:"kind" db:"kind" example:"regular"`                             // Тип кошелька (regular/bank)
	IsActive   bool       `json:"is_active" db:"is_active" example:"true"`                      // Статус активности кошелька
	CreatedAt  time.Time  `json:"created_at" db:"created_at" example:"2023-01-01T12:00:00Z"`    // Дата создания
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at" example:"2023-01-01T12:00:00Z"`    // Дата последнего обновления
	Username   string     `json:"username" db:"username" example:"user123"`                     // Имя пользователя владельца
}

// CreateWalletRequest запрос для создания нового кошелька
// @Description Данные, необходимые для создания нового блокчейн кошелька
type CreateWalletRequest struct {
	Kind     string `json:"kind" binding:"required" example:"regular"`     // Тип кошелька: regular или bank
	Username string `json:"username" binding:"required" example:"user123"` // Имя пользователя владельца кошелька
}

// WalletFilter фильтр для поиска кошельков
// @Description Параметры фильтрации при получении списка кошельков
type WalletFilter struct {
	Kind     string `json:"kind" example:"regular"`   // Фильтр по типу кошелька
	IsActive bool   `json:"is_active" example:"true"` // Фильтр по статусу активности
	Page     int    `json:"page" example:"0"`         // Номер страницы для пагинации
	Limit    int    `json:"limit" example:"10"`       // Количество записей на странице
}

// Pagination информация о пагинации
// @Description Метаданные пагинации для списочных запросов
type Pagination struct {
	Page  int `json:"page" example:"0"`    // Текущая страница
	Limit int `json:"limit" example:"10"`  // Количество записей на странице
	Total int `json:"total" example:"100"` // Общее количество записей
}

// WalletsResponse ответ со списком кошельков
// @Description Ответ API содержащий список кошельков и информацию о пагинации
type WalletsResponse struct {
	Wallets    []Wallet   `json:"wallets"`    // Список кошельков
	Pagination Pagination `json:"pagination"` // Информация о пагинации
}

// HTTPError структура для HTTP ошибок
// @Description Структура ошибки HTTP ответа
type HTTPError struct {
	Message string `json:"message" example:"Произошла ошибка"` // Сообщение об ошибке
}

// BalanceResponse ответ с балансом кошелька
// @Description Ответ содержащий информацию о балансе кошелька
type BalanceResponse struct {
	Address string      `json:"address" example:"TRX9sGPvkr7i3m1o..."` // Адрес кошелька
	Balance interface{} `json:"balance"`                               // Баланс кошелька (TRX и USDT)
}

// TransactionStatusResponse ответ со статусом транзакции
// @Description Ответ содержащий статус транзакции
type TransactionStatusResponse struct {
	TxID   string `json:"tx_id" example:"a1b2c3d4e5f6..."` // Идентификатор транзакции
	Status string `json:"status" example:"confirmed"`      // Статус транзакции
}

// TransactionsResponse ответ со списком транзакций
// @Description Ответ содержащий список транзакций с пагинацией
type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"` // Список транзакций
	Pagination   Pagination    `json:"pagination"`   // Информация о пагинации
}
