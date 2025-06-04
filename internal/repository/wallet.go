package repository

import (
	"blockchain-wallet/internal/domain"
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

type WalletRepository struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) *WalletRepository {
	// Проверяем подключение и схему
	var tableName string
	err := db.Get(&tableName, `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'wallet'`)
	if err != nil {
		log.Printf("Error checking table existence: %v", err)
	} else {
		log.Printf("Found table: %s", tableName)
	}

	return &WalletRepository{
		db: db,
	}
}

func (r *WalletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	// Логируем SQL запрос
	query := `
		INSERT INTO wallet (public_key, private_key, address, seed_phrase, kind, is_active, created_at, updated_at, username)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	log.Printf("Executing query: %s", query)
	log.Printf("With values: %+v", wallet)

	_, err := r.db.ExecContext(ctx, query,
		wallet.PublicKey,
		wallet.PrivateKey,
		wallet.Address,
		wallet.SeedPhrase,
		wallet.Kind,
		wallet.IsActive,
		wallet.CreatedAt,
		wallet.UpdatedAt,
		wallet.Username,
	)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}
	return nil
}

func (r *WalletRepository) FindAll(ctx context.Context, filter domain.WalletFilter) ([]domain.Wallet, domain.Pagination, error) {
	query := `
    SELECT 
      public_key,
      private_key,
      address,
      seed_phrase,
      kind,
      is_active,
      created_at,
      updated_at,
      username
    FROM wallet
    LIMIT $1 OFFSET $2
  `
	offset := (filter.Page - 1) * filter.Limit

	log.Printf("Executing query: %s with params: limit=%d, offset=%d", query, filter.Limit, offset)

	// Проверяем структуру таблицы
	var columns []string
	err := r.db.SelectContext(ctx, &columns, `
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = 'wallet' 
		ORDER BY ordinal_position
	`)
	if err != nil {
		log.Printf("Error getting table structure: %v", err)
	} else {
		log.Printf("Table columns: %v", columns)
	}

	var wallets []domain.Wallet
	err = r.db.SelectContext(ctx, &wallets, query,
		filter.Limit,
		offset,
	)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, domain.Pagination{}, err
	}

	// Получаем общее количество записей
	var total int
	countQuery := `SELECT COUNT(*) FROM wallet`
	err = r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		log.Printf("Error getting total count: %v", err)
		return nil, domain.Pagination{}, err
	}

	pagination := domain.Pagination{
		Page:  filter.Page,
		Limit: filter.Limit,
		Total: total,
	}

	log.Printf("Found %d wallets", len(wallets))
	return wallets, pagination, nil
}

func (r *WalletRepository) FindByAddress(ctx context.Context, address string) (*domain.Wallet, error) {
	query := `SELECT * FROM wallet WHERE address = $1`
	var wallet domain.Wallet
	err := r.db.GetContext(ctx, &wallet, query, address)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) Update(ctx context.Context, wallet *domain.Wallet) error {
	query := `
		UPDATE wallet
		SET private_key = $1,
			seed_phrase = $2,
			is_active = $3,
			updated_at = $4
		WHERE address = $5
	`
	_, err := r.db.ExecContext(ctx, query,
		wallet.PrivateKey,
		wallet.SeedPhrase,
		wallet.IsActive,
		wallet.UpdatedAt,
		wallet.Address,
	)
	if err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}
	return nil
}

func (r *WalletRepository) GetTransactionStatus(ctx context.Context, txID string) (string, error) {
	query := `SELECT status FROM transaction WHERE hash = $1`
	var status string
	err := r.db.GetContext(ctx, &status, query, txID)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction status: %w", err)
	}
	return status, nil
}

func (r *WalletRepository) GetTransactions(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, domain.Pagination, error) {
    query := `
        SELECT 
            hash,
            from_address,
            to_address,
            amount,
            status,
            confirmations,
            created_at,
            updated_at
        FROM transaction
        WHERE 
            ($1::text = '' OR from_address = $1::text)
            AND ($2::text = '' OR to_address = $2::text)
            AND ($3::transaction_status = '' OR status = $3::transaction_status)
        ORDER BY created_at DESC
        LIMIT $4 OFFSET $5
    `
    offset := (filter.Page - 1) * filter.Limit

    var transactions []domain.Transaction
    err := r.db.SelectContext(ctx, &transactions, query,
        filter.FromAddress,
        filter.ToAddress,
        filter.Status,
        filter.Limit,
        offset,
    )
    if err != nil {
        return nil, domain.Pagination{}, fmt.Errorf("failed to get transactions: %w", err)
    }

    // Получаем общее количество записей
    var total int
    countQuery := `
        SELECT COUNT(*) 
        FROM transaction 
        WHERE 
            ($1::text = '' OR from_address = $1::text)
            AND ($2::text = '' OR to_address = $2::text)
            AND ($3::transaction_status = '' OR status = $3::transaction_status)
    `
    err = r.db.GetContext(ctx, &total, countQuery,
        filter.FromAddress,
        filter.ToAddress,
        filter.Status,
    )
    if err != nil {
        return nil, domain.Pagination{}, fmt.Errorf("failed to get total count: %w", err)
    }

    pagination := domain.Pagination{
        Page:  filter.Page,
        Limit: filter.Limit,
        Total: total,
    }

    return transactions, pagination, nil
}
