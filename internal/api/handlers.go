// internal/api/handlers.go
package api

import (
	"blockchain-wallet/internal/config"
	"blockchain-wallet/internal/domain"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// Handler содержит все зависимости, необходимые для работы API
type Handler struct {
	wallet WalletService
}

// NewHandler создает новый экземпляр зависимостей
func NewHandler(cfg *config.Config, service WalletService) *Handler {
	// Здесь будет инициализация сервисов
	return &Handler{
		wallet: service,
	}
}

// getIntParam получает целочисленный параметр из запроса со значением по умолчанию
func getIntParam(c echo.Context, name string, defaultValue int) int {
	val := c.QueryParam(name)
	if val == "" {
		return defaultValue
	}
	if i, err := strconv.Atoi(val); err == nil {
		return i
	}
	return defaultValue
}

// CreateWallet создает новый кошелек
// @Summary Create a new wallet
// @Description Create a new wallet with the specified kind and username
// @ID create-wallet
// @Accept  json
// @Produce  json
// @Param   wallet  body  domain.CreateWalletRequest  true  "Wallet creation request"
// @Success 201 {object} domain.Wallet "Successfully created wallet"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /wallets [post]
func (h *Handler) CreateWallet(c echo.Context) error {
	var req domain.CreateWalletRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверное тело запроса")
	}

	wallet, err := h.wallet.CreateWallet(c.Request().Context(), req)
	if err != nil {
		if err.Error() == "username is required" || err.Error() == "kind is required" {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, wallet)
}

// GetWallets получает список кошельков с фильтрацией
// @Summary Get list of wallets
// @Description Get a list of wallets with optional filtering by kind and active status
// @ID get-wallets
// @Produce  json
// @Param   kind      query   string  false  "Filter by wallet kind (regular, bank)" Enums(regular, bank)
// @Param   is_active query   bool    false  "Filter by active status"
// @Param   page      query   int     false  "Page number" default(0)
// @Param   limit     query   int     false  "Number of wallets per page" default(10)
// @Success 200 {object} domain.WalletsResponse "List of wallets and pagination info"
// @Failure 500 {string} string "Internal server error"
// @Router /wallets [get]
func (h *Handler) GetWallets(c echo.Context) error {
	filter := domain.WalletFilter{
		Kind:     c.QueryParam("kind"),
		IsActive: strings.EqualFold("true", c.QueryParam("is_active")),
		Page:     getIntParam(c, "page", 0),
		Limit:    getIntParam(c, "limit", 10),
	}

	wallets, pagination, err := h.wallet.GetWallets(c.Request().Context(), filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, domain.WalletsResponse{
		Wallets:    wallets,
		Pagination: pagination,
	})
}

// SendTransaction создает новую транзакцию
// @Summary Send a transaction
// @Description Send a transaction from one address to another
// @ID send-transaction
// @Accept  json
// @Produce  json
// @Param   transaction  body  domain.CreateTransactionRequest  true  "Transaction details"
// @Success 201 {object} domain.Transaction "Successfully sent transaction"
// @Failure 400 {string} string "Bad request or invalid token type"
// @Failure 500 {string} string "Internal server error"
// @Router /transaction/send [post]
func (h *Handler) SendTransaction(c echo.Context) error {
	var req domain.CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверное тело запроса")
	}

	// Валидация типа токена
	if req.TokenType != "TRX" && req.TokenType != "USDT" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid token type. Must be either TRX or USDT")
	}

	transaction, err := h.wallet.SendTransaction(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, transaction)
}

// GetWalletTransactions получает список транзакций с фильтрацией
// @Summary Get transactions for a wallet
// @Description Get a list of transactions for a specific wallet address with pagination
// @ID get-wallet-transactions
// @Produce  json
// @Param   address  path  string  true  "Wallet address"
// @Param   page      query   int     false  "Page number" default(0)
// @Param   limit     query   int     false  "Number of transactions per page" default(10)
// @Success 200 {object} map[string]interface{} "List of transactions and pagination info"
// @Failure 400 {string} string "Address is required"
// @Failure 500 {string} string "Internal server error"
// @Router /{address}/transactions [get]
func (h *Handler) GetWalletTransactions(c echo.Context) error {
	address := c.Param("address")
	if address == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Address is required")
	}

	pg := &domain.Pagination{
		Limit: getIntParam(c, "limit", 10),
		Page:  getIntParam(c, "page", 0),
	}

	transactions, p, err := h.wallet.GetWalletTransactions(c.Request().Context(), address, pg)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactions": transactions,
		"pagination":   p,
	})
}

// GetBalance получает баланс кошелька
// @Summary Get wallet balance
// @Description Get the balance for a specific wallet address
// @ID get-wallet-balance
// @Produce  json
// @Param   address  path  string  true  "Wallet address"
// @Success 200 {object} map[string]interface{} "Wallet balance"
// @Failure 400 {string} string "Invalid address or address is required"
// @Failure 404 {string} string "Wallet not found"
// @Failure 500 {string} string "Internal server error"
// @Router /wallets/{address}/balance [get]
func (h *Handler) GetBalance(c echo.Context) error {
	address := c.Param("address")
	if address == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Address is required")
	}

	balance, err := h.wallet.GetBalance(c.Request().Context(), address)
	if err != nil {
		if errors.Is(err, domain.ErrWalletNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"address": address,
		"balance": balance,
	})
}

// GetTransactionStatus получает статус транзакции
// @Summary Get transaction status
// @Description Get the status of a specific transaction by its ID
// @ID get-transaction-status
// @Produce  json
// @Param   tx_id  path  string  true  "Transaction ID"
// @Success 200 {object} map[string]interface{} "Transaction status"
// @Failure 400 {string} string "Transaction ID is required"
// @Failure 500 {string} string "Internal server error"
// @Router /transactions/{tx_id}/status [get]
func (h *Handler) GetTransactionStatus(c echo.Context) error {
	txID := c.Param("tx_id")
	if txID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Transaction ID is required")
	}

	status, err := h.wallet.GetTransactionStatus(c.Request().Context(), txID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"tx_id":  txID,
		"status": status,
	})
}

// RegisterRoutes регистрирует все маршруты API
func RegisterRoutes(e *echo.Echo, handler *Handler) {
	v1 := e.Group("/api/v1")
	{
		v1.POST("/wallets", handler.CreateWallet)                           // Создание кошелька
		v1.GET("/wallets", handler.GetWallets)                              // Получение списка кошельков
		v1.GET("/wallets/:address/balance", handler.GetBalance)             // Получение баланса кошелька
		v1.POST("/transaction/send", handler.SendTransaction)               // Отправка транзакции
		v1.GET("/:address/transactions", handler.GetWalletTransactions)     // Получение списка транзакций
		v1.GET("/transactions/:tx_id/status", handler.GetTransactionStatus) // Получение статуса транзакции
	}
}
