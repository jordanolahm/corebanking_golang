package dto

type TransactionRequest struct {
	AccountID       string `json:"accountId"`
	OperationTypeID int    `json:"operationTypeId"`
	Amount          int64  `json:"amount"`
}

func NewTransactionRequest(accountID string, operationTypeID int, amount int64) TransactionRequest {
	return TransactionRequest{
		AccountID:       accountID,
		OperationTypeID: operationTypeID,
		Amount:          amount,
	}
}

func (transactionValue *TransactionRequest) GetAccountID() string {
	return transactionValue.AccountID
}

func (transactionValue *TransactionRequest) GetOperationTypeID() int {
	return transactionValue.OperationTypeID
}

func (transactionValue *TransactionRequest) GetAmount() int64 {
	return transactionValue.Amount
}

func (transactionValue *TransactionRequest) SetAccountID(accountID string) {
	transactionValue.AccountID = accountID
}

func (transactionValue *TransactionRequest) SetOperationTypeID(operationTypeID int) {
	transactionValue.OperationTypeID = operationTypeID
}

func (transactionValue *TransactionRequest) SetAmount(amount int64) {
	transactionValue.Amount = amount
}
