package repository

import (
	"corebanking/internal/domain"
	"sync"
)

type AccountRepository struct {
	accounts map[string]*domain.Account
	mu       sync.RWMutex
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{
		accounts: make(map[string]*domain.Account),
	}
}

func (r *AccountRepository) FindById(id string) (*domain.Account, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, exists := r.accounts[id]
	return account, exists
}

func (r *AccountRepository) Save(account *domain.Account) *domain.Account {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.accounts[account.ID] = account
	return account
}

func (r *AccountRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.accounts = make(map[string]*domain.Account)
}
