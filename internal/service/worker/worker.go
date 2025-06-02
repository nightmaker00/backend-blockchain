package worker

import (
	"blockchain-wallet/internal/api"
	"blockchain-wallet/internal/domain"
	"context"
	"sync"
	"time"
)

type TransactionWorker struct {
	walletService api.WalletService
	pendingTxs    map[string]*domain.Transaction
	mu            sync.RWMutex
	stopChan      chan struct{}
}

func NewTransactionWorker(walletService api.WalletService) *TransactionWorker {
	return &TransactionWorker{
		walletService: walletService,
		pendingTxs:    make(map[string]*domain.Transaction),
		stopChan:      make(chan struct{}),
	}
}

func (w *TransactionWorker) Init(ctx context.Context) error {
	// Получаем все pending транзакции из БД
	filter := domain.TransactionFilter{
		Status: "pending",
		Page:   1,
		Limit:  1000,
	}
	txs, _, err := w.walletService.GetTransactions(ctx, filter)
	if err != nil {
		return err
	}

	// Заполняем мапу
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, tx := range txs {
		w.pendingTxs[tx.Hash] = &tx
	}

	return nil
}

func (w *TransactionWorker) Run(ctx context.Context) {
	// Запускаем горутину для проверки новых транзакций
	go w.checkNewTransactions(ctx)
	// Запускаем горутину для проверки статуса существующих транзакций
	go w.checkTransactionStatus(ctx)
}

func (w *TransactionWorker) checkNewTransactions(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopChan:
			return
		case <-ticker.C:
			// Получаем все regular кошельки
			filter := domain.WalletFilter{
				Kind:  "regular",
				Page:  1,
				Limit: 1000,
			}
			wallets, _, err := w.walletService.GetWallets(ctx, filter)
			if err != nil {
				continue
			}

			// Для каждого кошелька проверяем новые транзакции
			for _, wallet := range wallets {
				txs, err := w.walletService.GetWalletTransactions(ctx, wallet.Address)
				if err != nil {
					continue
				}

				// Добавляем новые транзакции в мапу
				w.mu.Lock()
				for _, tx := range txs {
					if _, exists := w.pendingTxs[tx.Hash]; !exists {
						w.pendingTxs[tx.Hash] = &tx
					}
				}
				w.mu.Unlock()
			}
		}
	}
}

func (w *TransactionWorker) checkTransactionStatus(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopChan:
			return
		case <-ticker.C:
			w.mu.RLock()
			// Создаем копию мапы для итерации
			txs := make(map[string]*domain.Transaction, len(w.pendingTxs))
			for k, v := range w.pendingTxs {
				txs[k] = v
			}
			w.mu.RUnlock()

			// Проверяем статус каждой транзакции
			for hash, tx := range txs {
				status, err := w.walletService.GetTransactionStatus(ctx, hash)
				if err != nil {
					continue
				}

				// Удаляем транзакцию из мапы если:
				// 1. Она подтверждена 20 раз
				// 2. Она завершилась с ошибкой
				// 3. Она в процессе обработки (processing) - оставляем в мапе
				if status == "confirmed" && tx.Confirmations >= 20 || status == "failed" {
					w.mu.Lock()
					delete(w.pendingTxs, hash)
					w.mu.Unlock()
				}
			}
		}
	}
}

func (w *TransactionWorker) Shutdown(ctx context.Context) error {
	close(w.stopChan)
	return nil
}
