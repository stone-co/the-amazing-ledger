package vos

import "github.com/stone-co/the-amazing-ledger/app"

// TODO: melhorar nome (Linguagem ubiquoa)
type Path struct {
	Account AccountPath
	Credit  int64
	Debit   int64
}

type SyntheticReport struct {
	TotalCredit int64
	TotalDebit  int64
	Paths       []Path
}

func NewSyntheticReport(totalCredit, totalDebit int64, accounts []Path) (*SyntheticReport, error) {
	if accounts == nil || len(accounts) < 1 {
		return nil, app.ErrInvalidSyntheticReportStructure
	}

	return &SyntheticReport{
		TotalCredit: totalCredit,
		TotalDebit:  totalDebit,
		Paths:       accounts,
	}, nil
}
