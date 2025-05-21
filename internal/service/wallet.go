package service

import (
	"blockchain-wallet/internal/domain"
	"blockchain-wallet/pkg/blockchain/tron"
	"context"
	"time"
)

// walletService реализует интерфейс WalletService
type walletService struct {
	repo WalletRepository
	tron tron.Client
}

// NewWalletService создает новый сервис для работы с кошельками
func NewWalletService(repo WalletRepository, tron tron.Client) *walletService {
	return &walletService{
		repo: repo,
		tron: tron,
	}
}

// CreateWallet создает новый кошелек в сети Tron и сохраняет его в базе данных
func (s *walletService) CreateWallet(ctx context.Context, req domain.CreateWalletRequest) (*domain.Wallet, error) {
	// Создание кошелька в Tron
	address, err := s.tron.CreateWallet(ctx)
	if err != nil {
		return nil, err
	}

	// Сохранение в БД
	wallet := &domain.Wallet{
		Address:    address,
		WalletType: req.WalletType,
		Name:       req.Name,
		Status:     "active",
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

// GetWallets возвращает список кошельков с учетом фильтрации
func (s *walletService) GetWallets(ctx context.Context, filter domain.WalletFilter) ([]domain.Wallet, domain.Pagination, error) {
	return s.repo.FindAll(ctx, filter)
}