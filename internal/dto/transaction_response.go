package dto

import "time"

type TransactionResponse struct {
	TransactionID   int64     `json:"transactionId"`
	AccountID       string    `json:"accountId"`
	OperationTypeID int       `json:"operationTypeId"`
	Amount          int64     `json:"amount"`
	EventDate       time.Time `json:"eventDate"`
}

func NewTransactionResponse(transactionID int64, accountID string, operationTypeID int, amount int64, eventDate time.Time) TransactionResponse {
	return TransactionResponse{
		TransactionID:   transactionID,
		AccountID:       accountID,
		OperationTypeID: operationTypeID,
		Amount:          amount,
		EventDate:       eventDate,
	}
}

func (transactionValue *TransactionResponse) GetTransactionID() int64 {
	return transactionValue.TransactionID
}

func (transactionValue *TransactionResponse) GetAccountID() string {
	return transactionValue.AccountID
}

func (transactionValue *TransactionResponse) GetOperationTypeID() int {
	return transactionValue.OperationTypeID
}

func (transactionValue *TransactionResponse) GetAmount() int64 {
	return transactionValue.Amount
}

func (transactionValue *TransactionResponse) GetEventDate() time.Time {
	return transactionValue.EventDate
}

func (transactionValue *TransactionResponse) SetAmount(amount int64) {
	transactionValue.Amount = amount
}

func (transactionValue *TransactionResponse) SetEventDate(eventDate time.Time) {
	transactionValue.EventDate = eventDate
}
