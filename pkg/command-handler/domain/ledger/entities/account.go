package entities

import "time"

const (
	Asset     AccountType = "asset"
	Liability AccountType = "liability"
)

type AccountType string

type Account struct {
	ID        string // primary key
	Type      AccountType
	OwnerID   string
	Owner     string
	Name      string
	Metadata  []string
	Balance   int
	CreatedAt time.Time
	UpdatedAt *time.Time
}
