package entities

type AccountInfo struct {
	AccountID      string
	CurrentVersion Version
	TotalCredit    int
	TotalDebit     int
}

func NewAccountInfo(accountID string, version Version, totalCredit, totalDebit int) *AccountInfo {
	return &AccountInfo{
		AccountID:      accountID,
		CurrentVersion: version,
		TotalCredit:    totalCredit,
		TotalDebit:     totalDebit,
	}
}

func (a AccountInfo) Balance() int {
	return a.TotalCredit - a.TotalDebit
}
