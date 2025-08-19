package service

import (
	"corebanking/internal/domain"
	"corebanking/internal/dto"
	"corebanking/internal/repository"
	"fmt"
	"time"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	accountRepo     *repository.AccountRepository
}

func NewTransactionService(trRepo *repository.TransactionRepository, acRepo *repository.AccountRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: trRepo,
		accountRepo:     acRepo,
	}
}

func (s *TransactionService) CreateTransaction(req *dto.TransactionRequest) (*dto.TransactionResponse, error) {
	account, exists := s.accountRepo.FindById(req.AccountID)
	if !exists {
		return nil, fmt.Errorf("account not found")
	}

	amount := s.normalizeAmount(req.OperationTypeID, req.Amount)
	available := account.Balance + account.OverdraftLimit

	if amount < 0 && (available+amount) < 0 {
		return nil, fmt.Errorf("insufficient funds for transaction")
	}

	account.Balance += amount
	s.accountRepo.Save(account)

	transaction := &domain.Transaction{
		TransactionID:   domain.NextTransactionID(),
		AccountID:       req.AccountID,
		OperationTypeID: req.OperationTypeID,
		Amount:          amount,
		EventDate:       time.Now(),
	}

	s.transactionRepo.Save(transaction)

	return &dto.TransactionResponse{
		TransactionID:   transaction.TransactionID,
		AccountID:       transaction.AccountID,
		OperationTypeID: transaction.OperationTypeID,
		Amount:          transaction.Amount,
		EventDate:       transaction.EventDate,
	}, nil
}

func (s *TransactionService) GetTransactionByID(transactionID int64) (*dto.TransactionResponse, error) {
	transaction := s.transactionRepo.FindByID(transactionID)
	if transaction == nil {
		return nil, fmt.Errorf("transaction not found")
	}

	return &dto.TransactionResponse{
		TransactionID:   transaction.TransactionID,
		AccountID:       transaction.AccountID,
		OperationTypeID: transaction.OperationTypeID,
		Amount:          transaction.Amount,
		EventDate:       transaction.EventDate,
	}, nil
}

func (s *TransactionService) GetTransactionsToday() []*dto.TransactionResponse {
	today := time.Now()
	transactions := s.transactionRepo.FindAllTransactionOnDate(today)
	return s.mapTransactionsToResponse(transactions)
}

func (s *TransactionService) GetTransactionsInRange(begin, end time.Time) []*dto.TransactionResponse {
	transactions := s.transactionRepo.FindAllTransactionsBetweenDate(begin, end)
	return s.mapTransactionsToResponse(transactions)
}

func (s *TransactionService) GetTransactionsByType(operationTypeID int) []*dto.TransactionResponse {
	if operationTypeID < 1 || operationTypeID > 4 {
		return nil
	}
	transactions := s.transactionRepo.FindAllOperationTypeByID(operationTypeID)
	return s.mapTransactionsToResponse(transactions)
}

func (s *TransactionService) normalizeAmount(operationTypeID int, amount int64) int64 {
	switch operationTypeID {
	case 1, 2, 3:
		return -amount
	case 4:
		return amount
	default:
		panic("operation type doesn't exist")
	}
}

func (s *TransactionService) HandleTransaction(req *dto.EventRequest) (interface{}, error) {
	switch req.Type {
	case "deposit":
		return s.handleDeposit(req)
	case "withdraw":
		return s.handleWithdraw(req)
	case "transfer":
		return s.handleTransfer(req)
	default:
		return nil, fmt.Errorf("invalid event type")
	}
}

func (s *TransactionService) handleDeposit(req *dto.EventRequest) (map[string]*domain.Account, error) {
	// Recupera a conta do repositório
	account, exists := s.accountRepo.FindById(req.Destination)
	if !exists {
		account = domain.NewAccount(req.Destination, 0)
	}

	account.Balance += req.Amount
	s.accountRepo.Save(account)

	return map[string]*domain.Account{
		"destination": account,
	}, nil
}

func (s *TransactionService) handleWithdraw(req *dto.EventRequest) (map[string]*domain.Account, error) {
	account, exists := s.accountRepo.FindById(req.Origin)
	if !exists {
		return nil, fmt.Errorf("account not found")
	}

	available := account.Balance + account.OverdraftLimit
	if available < req.Amount {
		return nil, fmt.Errorf("insufficient funds, including overdraft")
	}

	account.Balance -= req.Amount
	s.accountRepo.Save(account)

	return map[string]*domain.Account{
		"origin": account,
	}, nil
}

func (s *TransactionService) handleTransfer(req *dto.EventRequest) (map[string]*domain.Account, error) {
	origin, exists := s.accountRepo.FindById(req.Origin)
	if !exists {
		return nil, fmt.Errorf("origin account not found")
	}

	destination, exists := s.accountRepo.FindById(req.Destination)
	if !exists {
		destination = domain.NewAccount(req.Destination, 0)
	}

	available := origin.Balance + origin.OverdraftLimit
	if available < req.Amount {
		return nil, fmt.Errorf("insufficient funds, including overdraft")
	}

	origin.Balance -= req.Amount
	destination.Balance += req.Amount

	s.accountRepo.Save(origin)
	s.accountRepo.Save(destination)

	return map[string]*domain.Account{
		"origin":      origin,
		"destination": destination,
	}, nil
}

func (s *TransactionService) mapTransactionsToResponse(transactions []*domain.Transaction) []*dto.TransactionResponse {
	result := make([]*dto.TransactionResponse, 0, len(transactions))
	for _, t := range transactions {
		result = append(result, &dto.TransactionResponse{
			TransactionID:   t.TransactionID,
			AccountID:       t.AccountID,
			OperationTypeID: t.OperationTypeID,
			Amount:          t.Amount,
			EventDate:       t.EventDate,
		})
	}
	return result
}
