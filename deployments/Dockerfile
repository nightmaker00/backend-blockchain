FROM golang:1.22-alpine

RUN apk add --no-cache gcc musl-dev git

RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN swag init -g cmd/app/main.go

RUN go build -o main ./cmd/app

EXPOSE 8080

CMD ["./main"] 