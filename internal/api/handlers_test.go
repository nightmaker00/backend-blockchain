package api

import (
	"blockchain-wallet/internal/domain"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) CreateWallet(ctx context.Context, req domain.CreateWalletRequest) (*domain.Wallet, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *MockWalletService) GetWallets(ctx context.Context, filter domain.WalletFilter) ([]domain.Wallet, domain.Pagination, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]domain.Wallet), args.Get(1).(domain.Pagination), args.Error(2)
}

func TestWalletHandler_CreateWallet(t *testing.T) {
	// Arrange
	e := echo.New()
	mockWalletService := new(MockWalletService)
	handler := NewWalletHandler(mockWalletService)

	reqBody := domain.CreateWalletRequest{
		WalletType: "tron",
		Name:       "Test Wallet",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/wallets", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedWallet := &domain.Wallet{
		Address:    "test_address",
		WalletType: "tron",
		Name:       "Test Wallet",
		Status:     "active",
	}

	mockWalletService.On("CreateWallet", mock.Anything, reqBody).Return(expectedWallet, nil)

	// Act
	err := handler.CreateWallet(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response domain.Wallet
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedWallet.Address, response.Address)
	assert.Equal(t, expectedWallet.WalletType, response.WalletType)
	assert.Equal(t, expectedWallet.Name, response.Name)

	mockWalletService.AssertExpectations(t)
}

func TestWalletHandler_GetWallets(t *testing.T) {
	// Arrange
	e := echo.New()
	mockWalletService := new(MockWalletService)
	handler := NewWalletHandler(mockWalletService)

	req := httptest.NewRequest(http.MethodGet, "/wallets?wallet_type=tron&status=active&page=1&limit=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedWallets := []domain.Wallet{
		{
			Address:    "test_address_1",
			WalletType: "tron",
			Name:       "Test Wallet 1",
			Status:     "active",
		},
		{
			Address:    "test_address_2",
			WalletType: "tron",
			Name:       "Test Wallet 2",
			Status:     "active",
		},
	}

	expectedPagination := domain.Pagination{
		Page:  1,
		Limit: 10,
		Total: 2,
	}

	mockWalletService.On("GetWallets", mock.Anything, domain.WalletFilter{
		WalletType: "tron",
		Status:     "active",
		Page:       1,
		Limit:      10,
	}).Return(expectedWallets, expectedPagination, nil)

	// Act
	err := handler.GetWallets(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response domain.WalletsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedWallets, response.Wallets)
	assert.Equal(t, expectedPagination, response.Pagination)

	mockWalletService.AssertExpectations(t)
}
