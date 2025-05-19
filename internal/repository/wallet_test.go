package repository

import (
	"blockchain-wallet/internal/domain"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestWalletRepository_Create(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewWalletRepository(db)
	ctx := context.Background()

	wallet := &domain.Wallet{
		Address:    "test_address",
		WalletType: "tron",
		Name:       "Test Wallet",
		Status:     "active",
		CreatedAt:  time.Now(),
	}

	// Act
	err := repo.Create(ctx, wallet)

	// Assert
	assert.NoError(t, err)

	// Verify
	var savedWallet domain.Wallet
	err = db.First(&savedWallet, "address = ?", wallet.Address).Error
	assert.NoError(t, err)
	assert.Equal(t, wallet.Address, savedWallet.Address)
	assert.Equal(t, wallet.WalletType, savedWallet.WalletType)
	assert.Equal(t, wallet.Name, savedWallet.Name)
	assert.Equal(t, wallet.Status, savedWallet.Status)
}

func TestWalletRepository_FindAll(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	repo := NewWalletRepository(db)
	ctx := context.Background()

	// Create test data
	wallets := []domain.Wallet{
		{
			Address:    "test_address_1",
			WalletType: "tron",
			Name:       "Test Wallet 1",
			Status:     "active",
			CreatedAt:  time.Now(),
		},
		{
			Address:    "test_address_2",
			WalletType: "tron",
			Name:       "Test Wallet 2",
			Status:     "active",
			CreatedAt:  time.Now(),
		},
	}

	for _, w := range wallets {
		err := repo.Create(ctx, &w)
		require.NoError(t, err)
	}

	filter := domain.WalletFilter{
		WalletType: "tron",
		Status:     "active",
		Page:       1,
		Limit:      10,
	}

	// Act
	foundWallets, pagination, err := repo.FindAll(ctx, filter)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, foundWallets, 2)
	assert.Equal(t, int64(2), pagination.Total)
	assert.Equal(t, 1, pagination.Page)
	assert.Equal(t, 10, pagination.Limit)
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&domain.Wallet{})
	require.NoError(t, err)

	return db
}
