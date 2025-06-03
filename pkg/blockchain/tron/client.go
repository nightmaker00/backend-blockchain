package tron

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"crypto/sha256"

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

func (t *TronClient) GetBalance(ctx context.Context, address string) (*WalletBalance, error) {
	url := fmt.Sprintf("https://apilist.tronscanapi.com/api/account/tokens?address=%s&start=0&limit=50&hidden=0&show=0&sortType=0&sortBy=0", address)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tronscanResp TronScanTokensResponse

	if err := json.NewDecoder(resp.Body).Decode(&tronscanResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	balance := &WalletBalance{
		TRXBalance: 0,
		USDTBalance: 0,
	}
	
	for _, token := range tronscanResp.Data {
		switch token.TokenID {
			case "_": // TRX token
			  if token.Quantity.Float64() > 0 {
				balance.TRXBalance = token.Quantity.Float64()
			  }
			case USDTContractAddress: // USDT token
			  if token.Quantity.Float64() > 0 {
				balance.USDTBalance = token.Quantity.Float64()
			  }
		}
	}
	return balance, nil	
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
	httpReq.Header.Set(os.Getenv("TRON_NODE_API_KEY"), t.apiKey)

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

	hash := sha256.Sum256(rawData)
	signature, err := crypto.Sign(hash[:], ecdsaPrivateKey)

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
	broadcastHttpReq.Header.Set(os.Getenv("TRON_NODE_API_KEY"), t.apiKey)

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
