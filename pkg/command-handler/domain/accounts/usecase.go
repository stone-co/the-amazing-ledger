package accounts

type UseCase interface {
	CreateAccount(AccountInput) error
	GetAccount(string) (Account, error)
	UpdateBalance(string, int) error
}

type AccountInput struct {
	Type     string
	OwnerID  string
	Owner    string
	Name     string
	Metadata []string
}

type Account struct {
	ID       string   `json:"id"`
	OwnerID  string   `json:"owner_id"`
	Type     string   `json:"type"`
	Balance  int      `json:"balance"`
	Owner    string   `json:"owner"`
	Name     string   `json:"name"`
	Metadata []string `json:"metadata"`
}
