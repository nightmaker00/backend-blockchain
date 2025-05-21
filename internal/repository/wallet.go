package repository

import (
	"blockchain-wallet/internal/domain"
	"context"

	"gorm.io/gorm"
)

// walletRepository реализует интерфейс WalletRepository
type walletRepository struct {
	db *gorm.DB
}

// NewWalletRepository создает новый репозиторий кошельков
func NewWalletRepository(db *gorm.DB) *walletRepository {
	return &walletRepository{db: db}
}

// Create сохраняет новый кошелек в базу данных
func (r *walletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}

// FindAll возвращает список кошельков с учетом фильтрации и пагинации
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

// FindByAddress находит кошелек по его адресу
func (r *walletRepository) FindByAddress(ctx context.Context, address string) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.db.WithContext(ctx).Where("address = ?", address).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// Update обновляет информацию о кошельке
func (r *walletRepository) Update(ctx context.Context, wallet *domain.Wallet) error {
	return r.db.WithContext(ctx).Save(wallet).Error
}
