package vos

type QueryBalance struct {
	Query   Account
	Balance int
}

func NewQueryBalance(query Account, balance int) QueryBalance {
	return QueryBalance{
		Query:   query,
		Balance: balance,
	}
}
