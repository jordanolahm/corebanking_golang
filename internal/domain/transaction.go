package domain

import (
	"sync/atomic"
	"time"
)

var transactionCounter int64 = 0

type Transaction struct {
	TransactionID   int64     `json:"transactionId"`
	AccountID       string    `json:"accountId"`
	OperationTypeID int       `json:"operationTypeId"`
	Amount          int64     `json:"amount"`
	EventDate       time.Time `json:"eventDate"`
}

func NewTransaction(accountId string, operationTypeId int, amount int64) *Transaction {
	return &Transaction{
		TransactionID:   NextTransactionID(),
		AccountID:       accountId,
		OperationTypeID: operationTypeId,
		Amount:          amount,
		EventDate:       time.Now(),
	}
}

func (id *Transaction) GetTransactionID() int64 {
	return id.TransactionID
}

func (accID *Transaction) GetAccountID() string {
	return accID.AccountID
}

func (typeID *Transaction) GetOperationTypeID() int {
	return typeID.OperationTypeID
}

func (amount *Transaction) GetAmount() int64 {
	return amount.Amount
}

func (date *Transaction) GetEventDate() time.Time {
	return date.EventDate
}

func NextTransactionID() int64 {
	return atomic.AddInt64(&transactionCounter, 1)
}
