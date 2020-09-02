package ledger

type AccountsUseCase interface {
	CreateAccount(AccountInput) error
	GetAccount(string) (Account, error)
	UpdateBalance(string, int) error
}

type AccountInput struct {
	Type     string   `json:"type"`
	OwnerID  string   `json:"owner_id"`
	Owner    string   `json:"owner"`
	Name     string   `json:"name"`
	Metadata []string `json:"metadata"`
}

type Account struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	OwnerID  string   `json:"owner_id"`
	Owner    string   `json:"owner"`
	Name     string   `json:"name"`
	Metadata []string `json:"metadata"`
	Balance  int      `json:"balance"`
}
