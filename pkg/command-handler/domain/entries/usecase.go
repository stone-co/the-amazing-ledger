package entries

type UseCase interface {
	CreateTransaction([]EntryInput) error
}

type EntryInput struct {
	AccountType     string   `json:"account_type"`
	AccountOwnerID  string   `json:"account_owner_id"`
	AccountOwner    string   `json:"account_owner"`
	AccountName     string   `json:"account_name"`
	AccountMetadata []string `json:"account_metadata"`
	RequestID       string   `json:"request_id"`
	Amount          int      `json:"amount"`
}
