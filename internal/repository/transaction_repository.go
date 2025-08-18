package repository

import (
	"corebanking/internal/domain"
	"sync"
	"time"
)

type TransactionRepository struct {
	mu           sync.RWMutex
	transactions []*domain.Transaction
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{
		transactions: make([]*domain.Transaction, 0),
	}
}

// Save adiciona uma transação à lista
func (r *TransactionRepository) Save(transaction *domain.Transaction) *domain.Transaction {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.transactions = append(r.transactions, transaction)
	return transaction
}

// FindByID retorna uma transação pelo ID
func (r *TransactionRepository) FindByID(transactionID int64) *domain.Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.transactions {
		if t.TransactionID == transactionID {
			return t
		}
	}
	return nil
}

// FindAllOperationTypeByID retorna todas as transações de um tipo específico
func (r *TransactionRepository) FindAllOperationTypeByID(operationTypeID int) []*domain.Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Transaction, 0)
	for _, t := range r.transactions {
		if t.OperationTypeID == operationTypeID {
			result = append(result, t)
		}
	}
	return result
}

// FindAllTransactionsBetweenDate retorna transações dentro de um intervalo
func (r *TransactionRepository) FindAllTransactionsBetweenDate(begin, end time.Time) []*domain.Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Transaction, 0)
	for _, t := range r.transactions {
		if (t.EventDate.Equal(begin) || t.EventDate.After(begin)) &&
			(t.EventDate.Equal(end) || t.EventDate.Before(end)) {
			result = append(result, t)
		}
	}
	return result
}

// FindAllTransactionOnDate retorna transações de um dia específico
func (r *TransactionRepository) FindAllTransactionOnDate(date time.Time) []*domain.Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Transaction, 0)
	for _, t := range r.transactions {
		if t.EventDate.Year() == date.Year() &&
			t.EventDate.Month() == date.Month() &&
			t.EventDate.Day() == date.Day() {
			result = append(result, t)
		}
	}
	return result
}

// FindAll retorna todas as transações
func (r *TransactionRepository) FindAll() []*domain.Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.transactions
}

// Reset limpa todas as transações
func (r *TransactionRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.transactions = make([]*domain.Transaction, 0)
}
