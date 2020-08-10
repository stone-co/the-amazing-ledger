package transactions

type UseCase interface {
	CreateOperation(OperationInput) error
}

type TransactionInput struct {
	AccountType     string
	AccountOwnerID  string
	AccountOwner    string
	AccountName     string
	AccountMetadata []string
	RequestID       string
	TransactionType string
	Amount          int
}

type OperationInput struct {
	Operation []TransactionInput
}

// type Transaction struct {
// 	ID       string   `json:"id"`
// 	OwnerID  string   `json:"owner_id"`
// 	Type     string   `json:"type"`
// 	Balance  int      `json:"balance"`
// 	Owner    string   `json:"owner"`
// 	Name     string   `json:"name"`
// 	Metadata []string `json:"metadata"`
// }
