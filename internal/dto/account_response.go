package dto

type AccountResponse struct {
	AccountID      string `json:"accountId"`
	DocumentNumber string `json:"documentNumber"`
}

func NewAccountResponse(accountID, documentNumber string) AccountResponse {
	return AccountResponse{
		AccountID:      accountID,
		DocumentNumber: documentNumber,
	}
}
