package dto

type BalanceResponse struct {
	Balance int64 `json:"balance"`
}

func NewBalanceResponse(balance int64) BalanceResponse {
	return BalanceResponse{
		Balance: balance,
	}
}

func (b *BalanceResponse) SetBalance(balance int64) {
	b.Balance = balance
}

func (b *BalanceResponse) GetBalance() int64 {
	return b.Balance
}
