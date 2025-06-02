package tron

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
)

type TronClient struct {
	httpClient *http.Client
	apiURL     string
	apiKey     string
}

func NewClient(httpClient *http.Client, apiNodeKey, apiNodeURL string) *TronClient {
	return &TronClient{
		httpClient: httpClient,
		apiURL:     apiNodeURL,
		apiKey:     apiNodeKey,
	}
}

func (t *TronClient) GetBalance(ctx context.Context, address string) (float64, error) {
	// Создаем запрос на получение баланса
	req := tronRequest{
		JSONRPC: "2.0",
		Method:  "wallet/getaccount",
		Params: map[string]interface{}{
			"address": address,
		},
		ID: 1,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Создаем HTTP запрос
	httpReq, err := http.NewRequest("POST", t.apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("TRON-PRO-API-KEY", t.apiKey)

	// Отправляем запрос
	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return 0, fmt.Errorf("server error: received status code %d", resp.StatusCode)
	}

	var tronResp tronResponse
	if err := json.NewDecoder(resp.Body).Decode(&tronResp); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if tronResp.Error != nil {
		return 0, fmt.Errorf("tron api error: %s", tronResp.Error.Message)
	}

	// Декодируем ответ
	var account struct {
		Balance int64 `json:"balance"`
	}
	if err := json.Unmarshal(tronResp.Result, &account); err != nil {
		return 0, fmt.Errorf("failed to unmarshal account: %w", err)
	}

	// Конвертируем из SUN в TRX (1 TRX = 1,000,000 SUN)
	return float64(account.Balance) / 1_000_000, nil
}

func (t *TronClient) SendTransaction(ctx context.Context, fromAddress, toAddress string, amount float64, privateKey string) (*Transaction, error) {
	// Конвертируем сумму в SUN (1 TRX = 1,000,000 SUN)
	amountInSun := int64(amount * 1_000_000)

	// Создаем запрос на создание транзакции
	req := tronRequest{
		JSONRPC: "2.0",
		Method:  "wallet/createtransaction",
		Params: map[string]interface{}{
			"owner_address": fromAddress,
			"to_address":    toAddress,
			"amount":        amountInSun,
		},
		ID: 1,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Создаем HTTP запрос с API ключом
	httpReq, err := http.NewRequest("POST", t.apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("TRON-PRO-API-KEY", t.apiKey)

	// Отправляем запрос
	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("server error: received status code %d", resp.StatusCode)
	}

	var tronResp tronResponse
	if err := json.NewDecoder(resp.Body).Decode(&tronResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if tronResp.Error != nil {
		return nil, fmt.Errorf("tron api error: %s", tronResp.Error.Message)
	}

	// Декодируем транзакцию
	var tx *Transaction
	if err := json.Unmarshal(tronResp.Result, &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	// Конвертируем приватный ключ
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key format: %w", err)
	}

	ecdsaPrivateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key: %w", err)
	}

	// Подписываем транзакцию
	rawData, err := hex.DecodeString(tx.RawDataHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode raw data: %w", err)
	}

	signature, err := crypto.Sign(rawData, ecdsaPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Добавляем подпись к транзакции
	tx.Signature = []string{hex.EncodeToString(signature)}

	// Обновляем broadcast запрос с API ключом
	broadcastReq := tronRequest{
		JSONRPC: "2.0",
		Method:  "wallet/broadcasttransaction",
		Params:  tx,
		ID:      2,
	}

	broadcastBody, err := json.Marshal(broadcastReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal broadcast request: %w", err)
	}

	broadcastHttpReq, err := http.NewRequest("POST", t.apiURL, bytes.NewBuffer(broadcastBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create broadcast request: %w", err)
	}
	broadcastHttpReq.Header.Set("Content-Type", "application/json")
	broadcastHttpReq.Header.Set("TRON-PRO-API-KEY", t.apiKey)

	broadcastResp, err := t.httpClient.Do(broadcastHttpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast transaction: %w", err)
	}
	defer broadcastResp.Body.Close()

	if broadcastResp.StatusCode >= 500 {
		return nil, fmt.Errorf("server error during broadcast: received status code %d", broadcastResp.StatusCode)
	}

	var broadcastResult tronResponse
	if err := json.NewDecoder(broadcastResp.Body).Decode(&broadcastResult); err != nil {
		return nil, fmt.Errorf("failed to decode broadcast response: %w", err)
	}

	if broadcastResult.Error != nil {
		return nil, fmt.Errorf("broadcast error: %s", broadcastResult.Error.Message)
	}

	return tx, nil
}
