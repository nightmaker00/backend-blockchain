// internal/repository/wallet.go
package repository

import (
	"blockchain-wallet/internal/domain"
	"context"

	"gorm.io/gorm"
)

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) domain.WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}

func (r *walletRepository) FindAll(ctx context.Context, filter domain.WalletFilter) ([]domain.Wallet, domain.Pagination, error) {
	var wallets []domain.Wallet
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Wallet{})

	if filter.WalletType != "" {
		query = query.Where("wallet_type = ?", filter.WalletType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, domain.Pagination{}, err
	}

	err = query.Offset((filter.Page - 1) * filter.Limit).
		Limit(filter.Limit).
		Find(&wallets).Error

	pagination := domain.Pagination{
		Page:  filter.Page,
		Limit: filter.Limit,
		Total: int(total),
	}

	return wallets, pagination, err
}

func (r *walletRepository) FindByAddress(ctx context.Context, address string) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.db.WithContext(ctx).Where("address = ?", address).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) Update(ctx context.Context, wallet *domain.Wallet) error {
	return r.db.WithContext(ctx).Save(wallet).Error
}
