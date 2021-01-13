package vos

type AccountHistory struct {
	Account        AccountName
	TotalCredit    int
	TotalDebit     int
	EntriesHistory []EntryHistory
}

func NewAccountHistory(account AccountName, entriesHistory ...EntryHistory) (AccountHistory, error) {
	totalCredit := 0
	totalDebit := 0
	for _, entryHistory := range entriesHistory {
		if entryHistory.Operation == CreditOperation {
			totalCredit += entryHistory.Amount
		} else {
			totalDebit += entryHistory.Amount
		}
	}

	return AccountHistory{
		Account:        account,
		TotalCredit:    totalCredit,
		TotalDebit:     totalDebit,
		EntriesHistory: entriesHistory,
	}, nil
}

// append func to add entryHistory at accountHistory
