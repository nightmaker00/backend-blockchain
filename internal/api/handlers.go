// internal/api/handlers.go
package api

import (
	"strconv"
    "net/http"
    "github.com/labstack/echo/v4"
    "blockchain-wallet/internal/domain"
    "blockchain-wallet/internal/service"
)

type Dependencies struct {
    WalletService    service.WalletService
    TransactionService service.TransactionService
}

func NewDependencies(cfg *config.Config) *Dependencies {
    // Здесь будет инициализация сервисов
    return &Dependencies{}
}

type TransactionHandler struct {
    transactionService service.TransactionService
}

func NewTransactionHandler(ts service.TransactionService) *TransactionHandler {
    return &TransactionHandler{transactionService: ts}
}

// WalletHandler обрабатывает запросы, связанные с кошельками
type WalletHandler struct {
    walletService service.WalletService
}

func NewWalletHandler(ws service.WalletService) *WalletHandler {
    return &WalletHandler{walletService: ws}
}

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

func (h *WalletHandler) CreateWallet(c echo.Context) error {
    var req domain.CreateWalletRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
    }

    wallet, err := h.walletService.CreateWallet(c.Request().Context(), req)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusCreated, wallet)
}

func (h *WalletHandler) GetWallets(c echo.Context) error {
    filter := domain.WalletFilter{
        WalletType: c.QueryParam("wallet_type"),
        Status:     c.QueryParam("status"),
        Page:       getIntParam(c, "page", 1),
        Limit:      getIntParam(c, "limit", 10),
    }

    wallets, pagination, err := h.walletService.GetWallets(c.Request().Context(), filter)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusOK, domain.WalletsResponse{
        Wallets:    wallets,
        Pagination: pagination,
    })
}

// Регистрация маршрутов
func RegisterRoutes(e *echo.Echo, deps *Dependencies) {
    walletHandler := NewWalletHandler(deps.WalletService)
    transactionHandler := NewTransactionHandler(deps.TransactionService)

    v1 := e.Group("/api/v1")
    {
        v1.POST("/wallets", walletHandler.CreateWallet)
        v1.GET("/wallets", walletHandler.GetWallets)
        // ... другие маршруты
    }
}