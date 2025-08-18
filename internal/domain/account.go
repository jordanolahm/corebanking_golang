package domain

type Account struct {
	ID             string `json:"id"`
	Balance        int64  `json:"balance"`
	OverdraftLimit int64  `json:"overdraft_limit"`
}

func NewAccount(id string, balance int64) *Account {
	return &Account{
		ID:             id,
		Balance:        balance,
		OverdraftLimit: 0,
	}
}

func (acc *Account) GetID() string {
	return acc.ID
}

func (acc *Account) GetBalance() int64 {
	return acc.Balance
}

func (acc *Account) GetOverdraftLimit() int64 {
	return acc.OverdraftLimit
}

func (a *Account) SetBalance(balance int64) {
	a.Balance = balance
}

func (a *Account) SetOverdraftLimit(limit int64) {
	a.OverdraftLimit = limit
}
