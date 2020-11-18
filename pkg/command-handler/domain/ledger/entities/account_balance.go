package entities

type AccountBalance struct {
	AccountName    AccountName
	CurrentVersion Version
	TotalCredit    int
	TotalDebit     int
}

func NewAccountBalance(accountName AccountName, version Version, totalCredit, totalDebit int) *AccountBalance {
	return &AccountBalance{
		AccountName:    accountName,
		CurrentVersion: version,
		TotalCredit:    totalCredit,
		TotalDebit:     totalDebit,
	}
}

func (a AccountBalance) Balance() int {
	return a.TotalCredit - a.TotalDebit
}
