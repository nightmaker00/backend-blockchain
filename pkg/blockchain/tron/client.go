package tron

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/crypto"
)

type TronClient struct {
	httpClient *http.Client
	apiNodeURL string
	apiNodeKey string
	scanURL    string
	scanKey    string
}

func NewClient(httpClient *http.Client, apiNodeKey, apiNodeURL, scanKey, scanURL string) *TronClient {
	if apiNodeURL != "" && !strings.HasPrefix(apiNodeURL, "http://") && !strings.HasPrefix(apiNodeURL, "https://") {
		apiNodeURL = "https://" + apiNodeURL
	}

	// if scanURL != "" && !strings.HasPrefix(scanURL, "http://") && !strings.HasPrefix(scanURL, "https://") {
	// 	scanURL = "https://" + scanURL
	// }

	return &TronClient{
		httpClient: httpClient,
		apiNodeURL: apiNodeURL,
		apiNodeKey: apiNodeKey,
		scanURL:    scanURL,
		scanKey:    scanKey,
	}
}

func (t *TronClient) GetBalance(ctx context.Context, address string) (*WalletBalance, error) {
	fmt.Printf("------------------here-------------------; address = %s\n", address)

	url := fmt.Sprintf("%s/account/tokens?address=%s&start=0&limit=4&hidden=0&show=0&sortType=0&sortBy=0", t.scanURL, address)

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
		TRXBalance:  0,
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
	amountInSun := int64(amount * 1_000_000)

	createTxReq := map[string]interface{}{
		"owner_address": fromAddress,
		"to_address":    toAddress,
		"amount":        amountInSun,
		"visible":       true,
	}

	reqBody, err := json.Marshal(createTxReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	createURL := t.apiNodeURL + "/wallet/createtransaction"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", createURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Add("accept", "application/json")
	httpReq.Header.Add("content-type", "application/json")
	// if t.apiNodeKey != "" {
	// 	httpReq.Header.Set("TRON-PRO-API-KEY", t.apiNodeKey)
	// }

	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	var tx Transaction
	if err := json.NewDecoder(resp.Body).Decode(&tx); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key format: %w", err)
	}

	ecdsaPrivateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key: %w", err)
	}

	rawData, err := hex.DecodeString(tx.RawDataHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode raw data: %w", err)
	}

	hash := sha256.Sum256(rawData)

	signature, err := crypto.Sign(hash[:], ecdsaPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	tx.Signature = []string{hex.EncodeToString(signature)}

	broadcastBody, err := json.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal broadcast request: %w", err)
	}

	broadcastURL := t.apiNodeURL + "/wallet/broadcasttransaction"
	broadcastHttpReq, err := http.NewRequestWithContext(ctx, "POST", broadcastURL, bytes.NewBuffer(broadcastBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create broadcast request: %w", err)
	}

	broadcastHttpReq.Header.Set("Content-Type", "application/json")

	if t.apiNodeKey != "" {
		broadcastHttpReq.Header.Set("TRON-PRO-API-KEY", t.apiNodeKey)
	}

	broadcastResp, err := t.httpClient.Do(broadcastHttpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast transaction: %w", err)
	}
	defer broadcastResp.Body.Close()

	if broadcastResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("broadcast failed with status code: %d", broadcastResp.StatusCode)
	}

	var broadcastResult map[string]interface{}
	if err := json.NewDecoder(broadcastResp.Body).Decode(&broadcastResult); err != nil {
		return nil, fmt.Errorf("failed to decode broadcast response: %w", err)
	}

	if result, ok := broadcastResult["result"].(bool); !ok || !result {
		if message, exists := broadcastResult["message"]; exists {
			return nil, fmt.Errorf("broadcast failed: %v", message)
		}
		return nil, fmt.Errorf("broadcast failed: unknown error")
	}

	return &tx, nil
}

// convertTronAddressToHex конвертирует TRON адрес из base58 в hex формат
func convertTronAddressToHex(address string) (string, error) {
	decoded := base58.Decode(address)

	if len(decoded) != 25 {
		return "", fmt.Errorf("invalid TRON address length")
	}

	// Убираем первый байт (0x41) и последние 4 байта (checksum)
	hexAddress := hex.EncodeToString(decoded[1:21])
	return hexAddress, nil
}

func encodeTransferParameters(toAddress string, amount *big.Int) (string, error) {
	addressHex, err := convertTronAddressToHex(toAddress)
	if err != nil {
		return "", fmt.Errorf("failed to convert address: %w", err)
	}

	paddedAddress := fmt.Sprintf("%064s", addressHex)

	amountHex := fmt.Sprintf("%064s", amount.Text(16))

	parameters := paddedAddress + amountHex
	return parameters, nil
}

func (t *TronClient) SendTRC20Token(ctx context.Context, fromAddress, toAddress string, amount float64, privateKey string, tokenType TokenType) (*TRC20Transaction, error) {
	tokenConfig, exists := TokenConfigs[tokenType]
	if !exists {
		return nil, fmt.Errorf("unsupported token type: %s", tokenType)
	}

	if tokenConfig.ContractAddress == "" {
		return nil, fmt.Errorf("token %s doesn't have contract address, use SendTransaction for TRX", tokenType)
	}

	amountInUnits := big.NewInt(int64(amount * float64(pow10(tokenConfig.Decimals, 10))))

	parameters, err := encodeTransferParameters(toAddress, amountInUnits)
	if err != nil {
		return nil, fmt.Errorf("failed to encode parameters: %w", err)
	}

	triggerReq := map[string]interface{}{
		"contract_address":  tokenConfig.ContractAddress,
		"function_selector": "transfer(address,uint256)",
		"parameter":         parameters,
		"fee_limit":         30000000, // 30 TRX fee limit
		"call_value":        0,
		"owner_address":     fromAddress,
	}

	reqBody, err := json.Marshal(triggerReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	triggerURL := t.apiNodeURL + "/wallet/triggersmartcontract"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", triggerURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if t.apiNodeKey != "" {
		httpReq.Header.Set("TRON-PRO-API-KEY", t.apiNodeKey)
	}

	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	var triggerResp struct {
		Result struct {
			Result bool `json:"result"`
		} `json:"result"`
		Transaction Transaction `json:"transaction"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&triggerResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !triggerResp.Result.Result {
		return nil, fmt.Errorf("failed to create TRC-20 transaction")
	}

	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	ecdsaPrivateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key: %w", err)
	}

	rawData, err := hex.DecodeString(triggerResp.Transaction.RawDataHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode raw data: %w", err)
	}

	hash := sha256.Sum256(rawData)
	signature, err := crypto.Sign(hash[:], ecdsaPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	trc20Tx := &TRC20Transaction{
		TxID:             triggerResp.Transaction.TxID,
		ContractAddress:  tokenConfig.ContractAddress,
		FunctionSelector: "transfer(address,uint256)",
		Parameter:        parameters,
		FeeLimit:         30000000,
		CallValue:        0,
		OwnerAddress:     fromAddress,
		RawDataHex:       triggerResp.Transaction.RawDataHex,
		Signature:        []string{hex.EncodeToString(signature)},
	}

	var txForBroadcast Transaction = triggerResp.Transaction
	txForBroadcast.Signature = []string{hex.EncodeToString(signature)}

	broadcastBody, err := json.Marshal(txForBroadcast)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal broadcast request: %w", err)
	}

	broadcastURL := t.apiNodeURL + "/wallet/broadcasttransaction"
	broadcastHttpReq, err := http.NewRequestWithContext(ctx, "POST", broadcastURL, bytes.NewBuffer(broadcastBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create broadcast request: %w", err)
	}

	broadcastHttpReq.Header.Set("Content-Type", "application/json")
	if t.apiNodeKey != "" {
		broadcastHttpReq.Header.Set("TRON-PRO-API-KEY", t.apiNodeKey)
	}

	broadcastResp, err := t.httpClient.Do(broadcastHttpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast transaction: %w", err)
	}
	defer broadcastResp.Body.Close()

	if broadcastResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("broadcast failed with status code: %d", broadcastResp.StatusCode)
	}

	var broadcastResult map[string]interface{}
	if err := json.NewDecoder(broadcastResp.Body).Decode(&broadcastResult); err != nil {
		return nil, fmt.Errorf("failed to decode broadcast response: %w", err)
	}

	if result, ok := broadcastResult["result"].(bool); !ok || !result {
		if message, exists := broadcastResult["message"]; exists {
			return nil, fmt.Errorf("broadcast failed: %v", message)
		}
		return nil, fmt.Errorf("broadcast failed: unknown error")
	}

	return trc20Tx, nil
}

// SendToken универсальный метод для отправки любых токенов
func (t *TronClient) SendToken(ctx context.Context, fromAddress, toAddress string, amount float64, privateKey string, tokenType TokenType) (interface{}, error) {
	switch tokenType {
	case TokenTypeTRX:
		return t.SendTransaction(ctx, fromAddress, toAddress, amount, privateKey)
	default:
		return t.SendTRC20Token(ctx, fromAddress, toAddress, amount, privateKey, tokenType)
	}
}

// pow10 возвращает n в степени k
func pow10(n, k int) int64 {
	result := int64(1)
	for i := 0; i < k; i++ {
		result *= int64(n)
	}
	return result
}
