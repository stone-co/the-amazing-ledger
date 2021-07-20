package vos

type AccountBalance struct {
	Account        Account
	CurrentVersion Version
	TotalCredit    int
	TotalDebit     int
}

func NewAccountBalance(account Account, version Version, totalCredit, totalDebit int) AccountBalance {
	return AccountBalance{
		Account:        account,
		CurrentVersion: version,
		TotalCredit:    totalCredit,
		TotalDebit:     totalDebit,
	}
}

func (a AccountBalance) Balance() int {
	return a.TotalCredit - a.TotalDebit
}
