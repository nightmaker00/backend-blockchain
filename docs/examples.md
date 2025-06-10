# Примеры использования Blockchain Wallet API

## Создание кошелька

### Запрос
```bash
curl -X POST http://localhost:8080/api/v1/wallets \
  -H "Content-Type: application/json" \
  -d '{
    "kind": "regular",
    "username": "user123"
  }'
```

### Ответ
```json
{
  "public_key": "0x1234567890abcdef...",
  "private_key": "0xabcdef1234567890...",
  "address": "TRX9sGPvkr7i3m1o...",
  "seed_phrase": "word1 word2 word3...",
  "kind": "regular",
  "is_active": true,
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z",
  "username": "user123"
}
```

## Получение списка кошельков

### Запрос
```bash
curl -X GET "http://localhost:8080/api/v1/wallets?kind=regular&is_active=true&page=0&limit=10"
```

### Ответ
```json
{
  "wallets": [
    {
      "public_key": "0x1234567890abcdef...",
      "private_key": "0xabcdef1234567890...",
      "address": "TRX9sGPvkr7i3m1o...",
      "seed_phrase": "word1 word2 word3...",
      "kind": "regular",
      "is_active": true,
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z",
      "username": "user123"
    }
  ],
  "pagination": {
    "page": 0,
    "limit": 10,
    "total": 100
  }
}
```

## Получение баланса кошелька

### Запрос
```bash
curl -X GET http://localhost:8080/api/v1/wallets/TRX9sGPvkr7i3m1o.../balance
```

### Ответ
```json
{
  "address": "TRX9sGPvkr7i3m1o...",
  "balance": {
    "trx": "1000.50",
    "usdt": "2500.75"
  }
}
```

## Отправка транзакции

### Запрос
```bash
curl -X POST http://localhost:8080/api/v1/transaction/send \
  -H "Content-Type: application/json" \
  -d '{
    "from_address": "TRX9sGPvkr7i3m1o...",
    "to_address": "TRX8sHQmkc6i4n2p...",
    "amount": 100.50,
    "token_type": "TRX"
  }'
```

### Ответ
```json
{
  "hash": "a1b2c3d4e5f6...",
  "from_address": "TRX9sGPvkr7i3m1o...",
  "to_address": "TRX8sHQmkc6i4n2p...",
  "amount": 100.50,
  "status": "pending",
  "confirmations": 0,
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z"
}
```

## Получение транзакций кошелька

### Запрос
```bash
curl -X GET "http://localhost:8080/api/v1/TRX9sGPvkr7i3m1o.../transactions?limit=10&page=0"
```

### Ответ
```json
{
  "transactions": [
    {
      "hash": "a1b2c3d4e5f6...",
      "from_address": "TRX9sGPvkr7i3m1o...",
      "to_address": "TRX8sHQmkc6i4n2p...",
      "amount": 100.50,
      "status": "confirmed",
      "confirmations": 12,
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z"
    }
  ],
  "pagination": {
    "page": 0,
    "limit": 10,
    "total": 50
  }
}
```

## Получение статуса транзакции

### Запрос
```bash
curl -X GET http://localhost:8080/api/v1/transactions/a1b2c3d4e5f6.../status
```

### Ответ
```json
{
  "tx_id": "a1b2c3d4e5f6...",
  "status": "confirmed"
}
```

## Коды ошибок

### 400 Bad Request
```json
{
  "message": "Некорректный запрос"
}
```

### 404 Not Found
```json
{
  "message": "Кошелек не найден"
}
```

### 500 Internal Server Error
```json
{
  "message": "Внутренняя ошибка сервера"
}
``` 