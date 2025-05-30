package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"blockchain-wallet/internal/domain"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
)

type walletService struct {
	tc   TronClient
	repo WalletRepository
}

func NewWalletService(tc TronClient, repo WalletRepository) *walletService {
	return &walletService{
		tc:   tc,
		repo: repo,
	}
}

func (s *walletService) CreateWallet(ctx context.Context, req domain.CreateWalletRequest) (*domain.Wallet, error) {
	// Генерируем энтропию (128 бит = 12 слов)
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, fmt.Errorf("failed to generate entropy: %w", err)
	}

	// Создаем мнемоническую фразу
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	// Получаем seed из мнемонической фразы
	seed := bip39.NewSeed(mnemonic, "")

	// Создаем мастер-ключ
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key: %w", err)
	}

	// Путь деривации для Tron: m/44'/195'/0'/0/0
	child, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		return nil, fmt.Errorf("failed to derive path: %w", err)
	}

	child, err = child.NewChildKey(bip32.FirstHardenedChild + 195)
	if err != nil {
		return nil, fmt.Errorf("failed to derive path: %w", err)
	}

	child, err = child.NewChildKey(bip32.FirstHardenedChild)
	if err != nil {
		return nil, fmt.Errorf("failed to derive path: %w", err)
	}

	child, err = child.NewChildKey(0)
	if err != nil {
		return nil, fmt.Errorf("failed to derive path: %w", err)
	}

	child, err = child.NewChildKey(0)
	if err != nil {
		return nil, fmt.Errorf("failed to derive path: %w", err)
	}

	// Получаем приватный ключ
	privateKey, err := crypto.ToECDSA(child.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key: %w", err)
	}

	// Получаем публичный ключ
	publicKey := privateKey.PublicKey
	publicKeyBytes := crypto.FromECDSAPub(&publicKey)

	// Получаем адрес Tron
	keccak256 := sha3.NewLegacyKeccak256()
	keccak256.Write(publicKeyBytes[1:])
	addressBytes := keccak256.Sum(nil)[12:]

	// Добавляем префикс Tron (0x41)
	tronAddress := append([]byte{0x41}, addressBytes...)

	// Вычисляем контрольную сумму
	firstSHA := sha256.Sum256(tronAddress)
	secondSHA := sha256.Sum256(firstSHA[:])
	checksum := secondSHA[:4]

	// Добавляем контрольную сумму к адресу
	finalAddress := append(tronAddress, checksum...)

	// Конвертируем в base64
	addressStr := base64.StdEncoding.EncodeToString(finalAddress)

	// Определяем тип кошелька
	kind := domain.WalletKindRegular
	if req.Kind != "" {
		kind = domain.WalletKind(req.Kind)
	}

	wallet := &domain.Wallet{
		PublicKey:  string(publicKeyBytes),
		PrivateKey: string(crypto.FromECDSA(privateKey)),
		Address:    addressStr,
		SeedPhrase: mnemonic,
		Kind:       kind,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Username:   req.Username,
	}

	if err := s.repo.Create(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to save wallet: %w", err)
	}
	

	return wallet, nil
}

func (s *walletService) GetBalance(ctx context.Context, address string) (float64, error) {
	return s.tc.GetBalance(ctx, address)
}

func (s *walletService) SendTransaction(ctx context.Context, req domain.CreateTransactionRequest) (*domain.Transaction, error) {
	// Получаем кошелек отправителя из базы данных
	wallet, err := s.repo.FindByAddress(ctx, req.FromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to find wallet: %w", err)
	}

	tx, err := s.tc.SendTransaction(ctx, req.FromAddress, req.ToAddress, req.Amount, wallet.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return domain.ToDomainTransaction(tx), nil
}

func (s *walletService) GetTransactionStatus(ctx context.Context, txID string) (string, error) {
	status, err := s.tc.GetTransactionStatus(ctx, txID)
	if err != nil {
		return "", fmt.Errorf("failed to get transaction status: %w", err)
	}
	return status, nil
}

func (w *walletService) GetWallets(ctx context.Context, filter domain.WalletFilter) ([]domain.Wallet, domain.Pagination, error) {
	return w.repo.FindAll(ctx, filter)
}
