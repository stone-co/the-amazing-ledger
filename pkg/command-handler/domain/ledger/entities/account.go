package entities

import "time"

const (
	Asset     AccountType = "asset"
	Liability AccountType = "liability"
)

type AccountType string

type Account struct {
	ID        string // primary key
	OwnerID   string
	Type      AccountType
	Balance   int
	Owner     string
	Name      string
	Metadata  []string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
