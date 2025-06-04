package tron

import (
	"encoding/json"
	"strconv"
)

type Transaction struct {
	Visible bool   `json:"visible"`
	TxID    string `json:"txID"`
	RawData struct {
		Contract []struct {
			Parameter struct {
				Value struct {
					Amount       int64  `json:"amount"`
					OwnerAddress string `json:"owner_address"`
					ToAddress    string `json:"to_address"`
				} `json:"value"`
				TypeURL string `json:"type_url"`
			} `json:"parameter"`
			Type string `json:"type"`
		} `json:"contract"`
		RefBlockBytes string `json:"ref_block_bytes"`
		RefBlockHash  string `json:"ref_block_hash"`
		Expiration    int64  `json:"expiration"`
		Timestamp     int64  `json:"timestamp"`
	} `json:"raw_data"`
	RawDataHex string   `json:"raw_data_hex"`
	Signature  []string `json:"signature"`
}

type FlexibleFloat float64

func (f *FlexibleFloat) UnmarshalJSON(data []byte) error {
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		*f = FlexibleFloat(num)
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str == "" {
		*f = 0
		return nil
	}

	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		*f = 0
		return nil
	}

	*f = FlexibleFloat(num)
	return nil
}

func (f FlexibleFloat) Float64() float64 {
	return float64(f)
}

type TronScanTokensResponse struct {
//	Total int                 `json:"total"`
	Data  []TronScanTokenInfo `json:"data"`
}

type TronScanTokenInfo struct {
	TokenID      string        `json:"tokenId"`
	Balance      string        `json:"balance"`
	TokenName    string        `json:"tokenName"`
	TokenAbbr    string        `json:"tokenAbbr"`
	TokenDecimal int           `json:"tokenDecimal"`
	TokenType    string        `json:"tokenType"`
	Quantity     FlexibleFloat `json:"quantity"`
	TokenCanShow int           `json:"tokenCanShow"`
}

type WalletBalance struct {
	TRXBalance  float64 `json:"trx_balance"`
	USDTBalance float64 `json:"usdt_balance"`
}

const (
	USDTContractAddress = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
)

type TokenType string

const (
  TokenTypeTRX  TokenType = "TRX"
  TokenTypeUSDT TokenType = "USDT"
)

type TokenConfig struct {
  ContractAddress string
  Decimals        int
}

var TokenConfigs = map[TokenType]TokenConfig{
  TokenTypeTRX: {
    ContractAddress: "", 
    Decimals:        6,  
  },
  TokenTypeUSDT: {
    ContractAddress: "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
    Decimals:        6,
  },
}

type TRC20Transaction struct {
	Visible          bool     `json:"visible"`
	TxID             string   `json:"txID"`
	ContractAddress  string   `json:"contract_address"`
	FunctionSelector string   `json:"function_selector"`
	Parameter        string   `json:"parameter"`
	FeeLimit         int64    `json:"fee_limit"`
	CallValue        int64    `json:"call_value"`
	OwnerAddress     string   `json:"owner_address"`
	RawDataHex       string   `json:"raw_data_hex"`
	Signature        []string `json:"signature"`
  }
  