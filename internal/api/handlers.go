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
// @Summary      Создание нового кошелька
// @Description  Создает новый блокчейн кошелек на TRON сети для указанного пользователя
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Param        request  body      domain.CreateWalletRequest  true  "Данные для создания кошелька"
// @Success      201      {object}  domain.Wallet               "Созданный кошелек"
// @Failure      400      {object}  domain.HTTPError            "Некорректный запрос"
// @Failure      500      {object}  domain.HTTPError            "Внутренняя ошибка сервера"
// @Router       /wallets [post]
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
// @Summary      Получение списка кошельков
// @Description  Возвращает список кошельков с возможностью фильтрации по типу и статусу активности
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Param        kind      query     string  false  "Тип кошелька (regular/bank)"
// @Param        is_active query     boolean false  "Статус активности кошелька"
// @Param        page      query     int     false  "Номер страницы" default(0)
// @Param        limit     query     int     false  "Количество записей на странице" default(10)
// @Success      200       {object}  domain.WalletsResponse      "Список кошельков с пагинацией"
// @Failure      500       {object}  domain.HTTPError            "Внутренняя ошибка сервера"
// @Router       /wallets [get]
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
// @Summary      Отправка транзакции
// @Description  Создает и отправляет транзакцию в блокчейн TRON. Поддерживает TRX и USDT токены
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        request  body      domain.CreateTransactionRequest  true  "Данные для создания транзакции"
// @Success      201      {object}  domain.Transaction                "Созданная транзакция"
// @Failure      400      {object}  domain.HTTPError                 "Некорректный запрос"
// @Failure      500      {object}  domain.HTTPError                 "Внутренняя ошибка сервера"
// @Router       /transaction/send [post]
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
// @Summary      Получение транзакций кошелька
// @Description  Возвращает список транзакций для указанного адреса кошелька с пагинацией
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        address   path      string  true   "Адрес кошелька"
// @Param        limit     query     int     false  "Количество записей на странице" default(10)
// @Param        page      query     int     false  "Номер страницы" default(0)
// @Success      200       {object}  domain.TransactionsResponse     "Список транзакций с пагинацией"
// @Failure      400       {object}  domain.HTTPError                "Некорректный запрос"
// @Failure      500       {object}  domain.HTTPError                "Внутренняя ошибка сервера"
// @Router       /{address}/transactions [get]
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

	return c.JSON(http.StatusOK, domain.TransactionsResponse{
		Transactions: transactions,
		Pagination:   p,
	})
}

// GetBalance получает баланс кошелька
// @Summary      Получение баланса кошелька
// @Description  Возвращает актуальный баланс TRX и USDT для указанного адреса кошелька
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Param        address   path      string  true  "Адрес кошелька в формате TRON"
// @Success      200       {object}  domain.BalanceResponse      "Баланс кошелька"
// @Failure      400       {object}  domain.HTTPError            "Некорректный запрос"
// @Failure      404       {object}  domain.HTTPError            "Кошелек не найден"
// @Failure      500       {object}  domain.HTTPError            "Внутренняя ошибка сервера"
// @Router       /wallets/{address}/balance [get]
func (h *Handler) GetBalance(c echo.Context) error {
	address := c.Param("address")
	if address == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Address is required")
	}

	balance, err := h.wallet.GetBalance(c.Request().Context(), address)
	if err != nil {
		if err.Error() == "wallet not found" {
			return echo.NewHTTPError(http.StatusNotFound, "Wallet not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, domain.BalanceResponse{
		Address: address,
		Balance: balance,
	})
}

// GetTransactionStatus получает статус транзакции
// @Summary      Получение статуса транзакции
// @Description  Возвращает текущий статус транзакции по её идентификатору
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        tx_id     path      string  true  "Идентификатор транзакции (hash)"
// @Success      200       {object}  domain.TransactionStatusResponse  "Статус транзакции"
// @Failure      400       {object}  domain.HTTPError                  "Некорректный запрос"
// @Failure      500       {object}  domain.HTTPError                  "Внутренняя ошибка сервера"
// @Router       /transactions/{tx_id}/status [get]
func (h *Handler) GetTransactionStatus(c echo.Context) error {
	txID := c.Param("tx_id")
	if txID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Transaction ID is required")
	}

	status, err := h.wallet.GetTransactionStatus(c.Request().Context(), txID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, domain.TransactionStatusResponse{
		TxID:   txID,
		Status: status,
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
