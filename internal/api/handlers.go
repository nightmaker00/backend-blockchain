// internal/api/handlers.go
package api

import (
	"blockchain-wallet/internal/config"
	"blockchain-wallet/internal/domain"
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
func (h *Handler) CreateWallet(c echo.Context) error {
	var req domain.CreateWalletRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверное тело запроса")
	}

	wallet, err := h.wallet.CreateWallet(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, wallet)
}

// GetWallets получает список кошельков с фильтрацией
func (h *Handler) GetWallets(c echo.Context) error {
	filter := domain.WalletFilter{
		Kind: 		c.QueryParam("kind"),
		IsActive:	strings.EqualFold("true", c.QueryParam("is_active")),
		Page:       getIntParam(c, "page", 1),
		Limit:      getIntParam(c, "limit", 10),
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

// CreateTransaction создает новую транзакцию
func (h *Handler) SendTransaction(c echo.Context) error {
	var req domain.CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверное тело запроса")
	}

	transaction, err := h.wallet.SendTransaction(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, transaction)
}

// GetTransactions получает список транзакций с фильтрацией
// func (h *Handler) GetTransactions(c echo.Context) error {
// 	filter := domain.TransactionFilter{
// 		FromAddress: c.QueryParam("from_address"),
// 		ToAddress:   c.QueryParam("to_address"),
// 		Status:      c.QueryParam("status"),
// 		Page:        getIntParam(c, "page", 1),
// 		Limit:       getIntParam(c, "limit", 10),
// 	}

// 	transactions, pagination, err := h.wallet.GetTransactions(c.Request().Context(), filter)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"transactions": transactions,
// 		"pagination":   pagination,
// 	})
// }

// RegisterRoutes регистрирует все маршруты API
func RegisterRoutes(e *echo.Echo, handler *Handler) {
	v1 := e.Group("/api/v1")
	{
		v1.POST("/wallets", handler.CreateWallet)
		v1.GET("/wallets", handler.GetWallets)
		v1.POST("/transaction/send", handler.SendTransaction)
		// v1.GET("/transactions", handler.GetTransactions)
	}
}
