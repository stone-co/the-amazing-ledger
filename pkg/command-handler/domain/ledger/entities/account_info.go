package entities

type AccountInfo struct {
	AccountPath    string
	CurrentVersion Version
	TotalCredit    int
	TotalDebit     int
}

func NewAccountInfo(accountPath string, version Version, totalCredit, totalDebit int) *AccountInfo {
	return &AccountInfo{
		AccountPath:    accountPath,
		CurrentVersion: version,
		TotalCredit:    totalCredit,
		TotalDebit:     totalDebit,
	}
}

func (a AccountInfo) Balance() int {
	return a.TotalCredit - a.TotalDebit
}
