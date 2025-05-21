# Используем официальный образ Go
FROM golang:1.21-alpine

# Устанавливаем необходимые зависимости
RUN apk add --no-cache gcc musl-dev

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/app

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"] 