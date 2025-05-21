package service

import (
	"blockchain-wallet/internal/domain"
	"blockchain-wallet/pkg/blockchain/tron"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWalletRepository реализует интерфейс repository.WalletRepository для тестирования
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) FindAll(ctx context.Context, filter domain.WalletFilter) ([]domain.Wallet, domain.Pagination, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]domain.Wallet), args.Get(1).(domain.Pagination), args.Error(2)
}

func (m *MockWalletRepository) FindByAddress(ctx context.Context, address string) (*domain.Wallet, error) {
	args := m.Called(ctx, address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *MockWalletRepository) Update(ctx context.Context, wallet *domain.Wallet) error {
	args := m.Called(ctx, wallet)
	return args.Error(0)
}

// MockTronClient реализует интерфейс tron.Client для тестирования
type MockTronClient struct {
	mock.Mock
}

func (m *MockTronClient) CreateWallet(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockTronClient) GetBalance(ctx context.Context, address string) (float64, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTronClient) GetTransactionStatus(ctx context.Context, txID string) (*tron.Transaction, error) {
	args := m.Called(ctx, txID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tron.Transaction), args.Error(1)
}

func (m *MockTronClient) SendTransaction(ctx context.Context, fromAddress, toAddress string, amount float64) (string, error) {
	args := m.Called(ctx, fromAddress, toAddress, amount)
	return args.String(0), args.Error(1)
}

// TestWalletService_CreateWallet тестирует создание кошелька
func TestWalletService_CreateWallet(t *testing.T) {
	// Подготовка
	mockRepo := new(MockWalletRepository)
	mockTron := new(MockTronClient)
	service := NewWalletService(mockRepo, mockTron)

	req := domain.CreateWalletRequest{
		WalletType: "tron",
		Name:       "Test Wallet",
	}

	expectedAddress := "test_address"
	mockTron.On("CreateWallet", mock.Anything).Return(expectedAddress, nil)

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(w *domain.Wallet) bool {
		return w.Address == expectedAddress &&
			w.WalletType == req.WalletType &&
			w.Name == req.Name &&
			w.Status == "active"
	})).Return(nil)

	// Действие
	wallet, err := service.CreateWallet(context.Background(), req)

	// Проверка
	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, expectedAddress, wallet.Address)
	assert.Equal(t, req.WalletType, wallet.WalletType)
	assert.Equal(t, req.Name, wallet.Name)
	assert.Equal(t, "active", wallet.Status)

	mockTron.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestWalletService_GetWallets тестирует получение списка кошельков
func TestWalletService_GetWallets(t *testing.T) {
	// Подготовка
	mockRepo := new(MockWalletRepository)
	mockTron := new(MockTronClient)
	service := NewWalletService(mockRepo, mockTron)

	filter := domain.WalletFilter{
		WalletType: "tron",
		Status:     "active",
		Page:       1,
		Limit:      10,
	}

	expectedWallets := []domain.Wallet{
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

	expectedPagination := domain.Pagination{
		Page:  1,
		Limit: 10,
		Total: 2,
	}

	mockRepo.On("FindAll", mock.Anything, filter).Return(expectedWallets, expectedPagination, nil)

	// Действие
	wallets, pagination, err := service.GetWallets(context.Background(), filter)

	// Проверка
	assert.NoError(t, err)
	assert.Equal(t, expectedWallets, wallets)
	assert.Equal(t, expectedPagination, pagination)

	mockRepo.AssertExpectations(t)
}

// TestWalletService_CreateWallet_Error тестирует обработку ошибки при создании кошелька
func TestWalletService_CreateWallet_Error(t *testing.T) {
	// Подготовка
	mockRepo := new(MockWalletRepository)
	mockTron := new(MockTronClient)
	service := NewWalletService(mockRepo, mockTron)

	req := domain.CreateWalletRequest{
		WalletType: "tron",
		Name:       "Test Wallet",
	}

	mockTron.On("CreateWallet", mock.Anything).Return("", assert.AnError)

	// Действие
	wallet, err := service.CreateWallet(context.Background(), req)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, wallet)
	assert.Equal(t, assert.AnError, err)

	mockTron.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "Create")
}
