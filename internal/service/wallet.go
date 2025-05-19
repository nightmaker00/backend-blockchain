// internal/service/wallet.go
package service

import (
	"blockchain-wallet/internal/domain"
	"blockchain-wallet/internal/tron"
	"context"
	"time"
)

type walletService struct {
	repo domain.WalletRepository
	tron tron.Client
}

func NewWalletService(repo domain.WalletRepository, tron tron.Client) domain.WalletService {
	return &walletService{
		repo: repo,
		tron: tron,
	}
}

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

func (s *walletService) GetWallets(ctx context.Context, filter domain.WalletFilter) ([]domain.Wallet, domain.Pagination, error) {
	return s.repo.FindAll(ctx, filter)
}
