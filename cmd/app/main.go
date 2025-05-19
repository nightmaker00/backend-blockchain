// cmd/api/main.go
package main

import (
    "log"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "blockchain-wallet/internal/api"
    "blockchain-wallet/internal/config"
)

func main() {
    // Загрузка конфигурации
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Инициализация сервера
    e := echo.New()
    
    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORS())

    // Инициализация зависимостей
    deps := api.NewDependencies(cfg)
    
    // Регистрация маршрутов
    api.RegisterRoutes(e, deps)

    // Запуск сервера
    log.Fatal(e.Start(cfg.Server.Address))
}