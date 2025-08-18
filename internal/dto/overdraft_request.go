package dto

type OverdraftRequest struct {
	AccountID string `json:"accountId"`
	Limit     int64  `json:"limit"`
}

func NewOverdraftRequest(accountID string, limit int64) OverdraftRequest {
	return OverdraftRequest{
		AccountID: accountID,
		Limit:     limit,
	}
}

func (overdraftValue *OverdraftRequest) GetAccountID() string {
	return overdraftValue.AccountID
}

func (overdraftValue *OverdraftRequest) GetLimit() int64 {
	return overdraftValue.Limit
}

func (overdraftValue *OverdraftRequest) SetAccountID(accountID string) {
	overdraftValue.AccountID = accountID
}

func (overdraftValue *OverdraftRequest) SetLimit(limit int64) {
	overdraftValue.Limit = limit
}
