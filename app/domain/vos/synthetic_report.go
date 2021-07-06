package vos

import "github.com/stone-co/the-amazing-ledger/app"

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

// synthetic account
func NewSyntheticReport(totalCredit, totalDebit int64, paths []Path) (*SyntheticReport, error) {
	if paths == nil || len(paths) < 1 {
		return nil, app.ErrInvalidSyntheticReportStructure
	}

	return &SyntheticReport{
		TotalCredit: totalCredit,
		TotalDebit:  totalDebit,
		Paths:       paths,
	}, nil
}
