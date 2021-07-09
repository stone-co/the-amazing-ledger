package vos

type QueryBalance struct {
	Query   AccountQuery
	Balance int
}

func NewQueryBalance(query AccountQuery, balance int) QueryBalance {
	return QueryBalance{
		Query:   query,
		Balance: balance,
	}
}
