package service

import (
	"corebanking/internal/domain"
	"corebanking/internal/dto"
	"corebanking/internal/repository"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type AccountService struct {
	accountRepo       *repository.AccountRepository
	documentToAccount map[string]string
	mu                sync.RWMutex
}

func NewAccountService(accountRepo *repository.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo:       accountRepo,
		documentToAccount: make(map[string]string),
	}
}

func (s *AccountService) CreateAccount(documentNumber string) (*dto.AccountResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.documentToAccount[documentNumber]; exists {
		return nil, fmt.Errorf("document already has an account")
	}

	accountID := uuid.New().String()
	account := &domain.Account{
		ID:      accountID,
		Balance: 0,
	}

	s.accountRepo.Save(account)
	s.documentToAccount[documentNumber] = accountID

	return &dto.AccountResponse{
		AccountID:      accountID,
		DocumentNumber: documentNumber,
	}, nil
}

func (s *AccountService) GetAccount(accountID string) (*dto.AccountResponse, error) {
	account, exists := s.accountRepo.FindById(accountID)
	if !exists {
		return nil, fmt.Errorf("account not found")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	document := "UNKNOWN"
	for doc, id := range s.documentToAccount {
		if id == accountID {
			document = doc
			break
		}
	}

	return &dto.AccountResponse{
		AccountID:      account.ID,
		DocumentNumber: document,
	}, nil
}

func (s *AccountService) GetBalance(accountID string) (*dto.BalanceResponse, error) {
	account, exists := s.accountRepo.FindById(accountID)
	if !exists {
		return nil, fmt.Errorf("account not found")
	}

	return &dto.BalanceResponse{
		Balance: account.Balance,
	}, nil
}

func (s *AccountService) ConfigOverdraft(accountID string, limit int64) error {
	account, exists := s.accountRepo.FindById(accountID)
	if !exists {
		return fmt.Errorf("account not found")
	}

	account.OverdraftLimit = limit
	s.accountRepo.Save(account)
	return nil
}

func (s *AccountService) Reset() {
	s.accountRepo.Reset()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.documentToAccount = make(map[string]string)
}
