# Makefile для Blockchain Wallet API

.PHONY: help build run test clean docs docs-serve docs-format docker-build docker-up docker-down docker-logs docker-swagger

# Переменные
BINARY_NAME=blockchain-wallet
DOCS_DIR=docs
MAIN_PATH=cmd/app/main.go
DOCKER_COMPOSE_FILE=deployments/docker-compose.yml

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Собрать приложение
	@echo "Сборка приложения..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)

run: ## Запустить приложение
	@echo "Запуск приложения..."
	go run $(MAIN_PATH)

test: ## Запустить тесты
	@echo "Запуск тестов..."
	go test -v ./...

clean: ## Очистить сборочные файлы
	@echo "Очистка файлов..."
	rm -f $(BINARY_NAME)
	rm -rf $(DOCS_DIR)

deps: ## Установить зависимости
	@echo "Установка зависимостей..."
	go mod download
	go mod tidy

docs: ## Генерировать Swagger документацию
	@echo "Генерация Swagger документации..."
	swag init -g $(MAIN_PATH)

docs-format: ## Форматировать Swagger аннотации
	@echo "Форматирование Swagger аннотаций..."
	swag fmt

docs-serve: docs ## Запустить сервер с документацией
	@echo "Запуск сервера с документацией..."
	@echo "Swagger UI будет доступен по адресу: http://localhost:8080/swagger/index.html"
	go run $(MAIN_PATH)

install-swag: ## Установить swag CLI
	@echo "Установка swag CLI..."
	go install github.com/swaggo/swag/cmd/swag@latest

dev: deps docs run ## Полная настройка для разработки

# Docker команды
docker-build: ## Собрать Docker образ
	@echo "Сборка Docker образа..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) build

docker-up: ## Запустить сервисы через Docker Compose
	@echo "Запуск сервисов через Docker Compose..."
	@echo "🚀 Запускаем blockchain wallet API..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo ""
	@echo "✅ Сервисы запущены!"
	@echo "📖 Swagger UI: http://localhost:8080/swagger/index.html"
	@echo "🔗 API базовый URL: http://localhost:8080/api/v1"
	@echo "📊 Логи приложения: make docker-logs"
	@echo "🛑 Остановить сервисы: make docker-down"

docker-down: ## Остановить сервисы Docker Compose
	@echo "Остановка сервисов Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

docker-logs: ## Показать логи всех сервисов
	@echo "Логи сервисов:"
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

docker-logs-api: ## Показать логи только API сервиса
	@echo "Логи API сервиса:"
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f api

docker-swagger: docker-up ## Запустить с Docker и открыть Swagger
	@echo "🔗 Swagger UI доступен по адресу: http://localhost:8080/swagger/index.html"
	@echo "Попробуйте открыть ссылку в браузере через несколько секунд..."

docker-restart: docker-down docker-up ## Перезапустить все сервисы

docker-clean: ## Очистить все Docker ресурсы проекта
	@echo "Очистка Docker ресурсов..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --rmi all --remove-orphans

.DEFAULT_GOAL := help 