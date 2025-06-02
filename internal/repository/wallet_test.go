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

// TestWalletRepository_Create_DuplicateUsername тестирует создание кошелька с дубликатом имени пользователя
func TestWalletRepository_Create_DuplicateUsername(t *testing.T) {
	// Подготовка
	db := setupTestDB(t)
	repo := NewWalletRepository(db)
	ctx := context.Background()

	wallet1 := &domain.Wallet{
		PublicKey:  "test_address_1",
		PrivateKey: "test_private_key_1",
		Address:    "test_address_1",
		SeedPhrase: "test_seed_phrase_1",
		Kind:       domain.WalletKindRegular,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Username:   "test_username",
	}

	wallet2 := &domain.Wallet{
		PublicKey:  "test_address_2",
		PrivateKey: "test_private_key_2",
		Address:    "test_address_2",
		SeedPhrase: "test_seed_phrase_2",
		Kind:       domain.WalletKindRegular,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Username:   "test_username",
	}

	// Действие
	err1 := repo.Create(ctx, wallet1)
	err2 := repo.Create(ctx, wallet2)

	// Проверка
	assert.NoError(t, err1)
	assert.Error(t, err2)
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
		{
			PublicKey:  "test_address_3",
			PrivateKey: "test_private_key_3",
			Address:    "test_address_3",
			SeedPhrase: "test_seed_phrase_3",
			Kind:       domain.WalletKindRegular,
			IsActive:   false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Username:   "test_username_3",
		},
	}

	for _, w := range wallets {
		err := repo.Create(ctx, &w)
		require.NoError(t, err)
	}

	tests := []struct {
		name     string
		filter   domain.WalletFilter
		expected int
	}{
		{
			name: "all wallets",
			filter: domain.WalletFilter{
				Page:  1,
				Limit: 10,
			},
			expected: 3,
		},
		{
			name: "active wallets",
			filter: domain.WalletFilter{
				IsActive: true,
				Page:     1,
				Limit:    10,
			},
			expected: 2,
		},
		{
			name: "regular wallets",
			filter: domain.WalletFilter{
				Kind:  "regular",
				Page:  1,
				Limit: 10,
			},
			expected: 3,
		},
		{
			name: "pagination",
			filter: domain.WalletFilter{
				Page:  1,
				Limit: 2,
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Действие
			foundWallets, pagination, err := repo.FindAll(ctx, tt.filter)

			// Проверка
			assert.NoError(t, err)
			assert.Len(t, foundWallets, tt.expected)
			assert.Equal(t, 3, pagination.Total)
			assert.Equal(t, tt.filter.Page, pagination.Page)
			assert.Equal(t, tt.filter.Limit, pagination.Limit)
		})
	}
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

	tests := []struct {
		name          string
		address       string
		shouldExist   bool
		expectedError bool
	}{
		{
			name:          "existing wallet",
			address:       wallet.Address,
			shouldExist:   true,
			expectedError: false,
		},
		{
			name:          "non-existing wallet",
			address:       "non_existing_address",
			shouldExist:   false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Действие
			foundWallet, err := repo.FindByAddress(ctx, tt.address)

			// Проверка
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, foundWallet)
			} else {
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
		})
	}
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
