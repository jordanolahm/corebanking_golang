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

// FindById retorna a conta e um boolean indicando se ela existe
func (r *AccountRepository) FindById(id string) (*domain.Account, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	account, exists := r.accounts[id]
	return account, exists
}

// Save insere ou atualiza uma conta
func (r *AccountRepository) Save(account *domain.Account) *domain.Account {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.accounts[account.ID] = account
	return account
}

// Reset limpa o repositório
func (r *AccountRepository) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.accounts = make(map[string]*domain.Account)
}
