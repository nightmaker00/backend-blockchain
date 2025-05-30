package repository

import (
	"blockchain-wallet/internal/domain"
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWalletRepository_Create тестирует создание кошелька в репозитории
func TestWalletRepository_Create(t *testing.T) {
	// Подготовка
	db := setupTestDB(t)
	repo := NewWalletRepository(db)
	ctx := context.Background()

	wallet := &domain.Wallet{
		PublicKey:  "test_address",
		PrivateKey: "test_private_key",
		Address:    "test_address",
		SeedPhrase: "test_seed_phrase",
		Kind:       domain.WalletKindRegular,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Username:   "test_username",
	}

	// Действие
	err := repo.Create(ctx, wallet)

	// Проверка
	assert.NoError(t, err)

	// Верификация
	var savedWallet domain.Wallet
	err = db.GetContext(ctx, &savedWallet, "SELECT * FROM wallet WHERE public_key = $1", wallet.PublicKey)
	assert.NoError(t, err)
	assert.Equal(t, wallet.PublicKey, savedWallet.PublicKey)
	assert.Equal(t, wallet.PrivateKey, savedWallet.PrivateKey)
	assert.Equal(t, wallet.SeedPhrase, savedWallet.SeedPhrase)
	assert.Equal(t, wallet.Kind, savedWallet.Kind)
	assert.Equal(t, wallet.IsActive, savedWallet.IsActive)
	assert.Equal(t, wallet.Username, savedWallet.Username)
}

// TestWalletRepository_FindAll тестирует получение списка кошельков с фильтрацией
func TestWalletRepository_FindAll(t *testing.T) {
	// Подготовка
	db := setupTestDB(t)
	repo := NewWalletRepository(db)
	ctx := context.Background()

	// Создание тестовых данных
	wallets := []domain.Wallet{
		{
			PublicKey:  "test_address_1",
			PrivateKey: "test_private_key_1",
			Address:    "test_address_1",
			SeedPhrase: "test_seed_phrase_1",
			Kind:       domain.WalletKindRegular,
			IsActive:   true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Username:   "test_username_1",
		},
		{
			PublicKey:  "test_address_2",
			PrivateKey: "test_private_key_2",
			Address:    "test_address_2",
			SeedPhrase: "test_seed_phrase_2",
			Kind:       domain.WalletKindRegular,
			IsActive:   true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Username:   "test_username_2",
		},
	}

	for _, w := range wallets {
		err := repo.Create(ctx, &w)
		require.NoError(t, err)
	}

	filter := domain.WalletFilter{
		Kind:     "regular",
		IsActive: true,
		Page:     1,
		Limit:    10,
	}

	// Действие
	foundWallets, pagination, err := repo.FindAll(ctx, filter)

	// Проверка
	assert.NoError(t, err)
	assert.Len(t, foundWallets, 2)
	assert.Equal(t, 2, pagination.Total)
	assert.Equal(t, 1, pagination.Page)
	assert.Equal(t, 10, pagination.Limit)
}

// TestWalletRepository_FindByAddress тестирует поиск кошелька по адресу
func TestWalletRepository_FindByAddress(t *testing.T) {
	// Подготовка
	db := setupTestDB(t)
	repo := NewWalletRepository(db)
	ctx := context.Background()

	wallet := &domain.Wallet{
		PublicKey:  "test_address",
		PrivateKey: "test_private_key",
		Address:    "test_address",
		SeedPhrase: "test_seed_phrase",
		Kind:       domain.WalletKindRegular,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Username:   "test_username",
	}

	err := repo.Create(ctx, wallet)
	require.NoError(t, err)

	// Действие
	foundWallet, err := repo.FindByAddress(ctx, wallet.Address)

	// Проверка
	assert.NoError(t, err)
	assert.NotNil(t, foundWallet)
	assert.Equal(t, wallet.Address, foundWallet.Address)
	assert.Equal(t, wallet.PublicKey, foundWallet.PublicKey)
	assert.Equal(t, wallet.PrivateKey, foundWallet.PrivateKey)
	assert.Equal(t, wallet.SeedPhrase, foundWallet.SeedPhrase)
	assert.Equal(t, wallet.Kind, foundWallet.Kind)
	assert.Equal(t, wallet.IsActive, foundWallet.IsActive)
	assert.Equal(t, wallet.Username, foundWallet.Username)
}

// setupTestDB создает тестовую базу данных
func setupTestDB(t *testing.T) *sqlx.DB {
	// Используем тестовую базу данных PostgreSQL
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=blockchain_wallet_test sslmode=disable")
	require.NoError(t, err)

	// Очищаем таблицу перед тестами
	_, err = db.Exec("DROP TABLE IF EXISTS wallet")
	require.NoError(t, err)

	// Создаем таблицу
	_, err = db.Exec(`
		CREATE TABLE wallet (
			public_key VARCHAR(255) PRIMARY KEY,
			private_key VARCHAR(255) NOT NULL,
			address VARCHAR(255) NOT NULL,
			seed_phrase TEXT NOT NULL,
			kind VARCHAR(50) NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			username VARCHAR(255) NOT NULL UNIQUE
		)
	`)
	require.NoError(t, err)

	return db
}
