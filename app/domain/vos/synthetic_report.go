package vos

import "github.com/stone-co/the-amazing-ledger/app"

type Path struct {
	Account string
	Credit  int
	Debit   int
}

type SyntheticReport struct {
	TotalCredit    int
	TotalDebit     int
	Paths          []Path
	CurrentVersion Version
}

// synthetic account
func NewSyntheticReport(totalCredit, totalDebit int, paths []Path, version Version) (*SyntheticReport, error) {
	if paths == nil || len(paths) < 1 {
		return nil, app.ErrInvalidSyntheticReportStructure
	}

	return &SyntheticReport{
		TotalCredit:    totalCredit,
		TotalDebit:     totalDebit,
		Paths:          paths,
		CurrentVersion: version,
	}, nil
}
