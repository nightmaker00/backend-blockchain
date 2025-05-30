package main

import (
	"blockchain-wallet/internal/api"
	"blockchain-wallet/internal/config"
	"blockchain-wallet/internal/repository"
	"blockchain-wallet/internal/service"
	tronlib "blockchain-wallet/pkg/blockchain/tron"
	pc "blockchain-wallet/pkg/db/postgres"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load("./deployments/.env.local"); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	dbClient, err := pc.NewPostgresDB(pc.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация репозитория
	walletRepo := repository.NewWalletRepository(dbClient)


	
	// Инициализация клиента Tron
	hc := &http.Client{}
	
	tcl := tronlib.NewClient(hc, os.Getenv("TRON_NODE_API_KEY"), os.Getenv("TRON_NODE_URL"))

	// Инициализация сервисов
	walletService := service.NewWalletService(tcl, walletRepo)

	// Инициализация сервера
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(middleware.CORS())

	// Инициализация обработчиков
	handler := api.NewHandler(cfg, walletService)

	// Регистрация маршрутов
	api.RegisterRoutes(e, handler)

	// Запуск сервера в горутине
	go func() {
		host := cfg.Server.Address + ":" + cfg.Server.Port
		log.Printf("Server starting on %s", host)

		if err := e.Start(host); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
