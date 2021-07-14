package vos

import (
	"time"

	"github.com/google/uuid"

	"github.com/stone-co/the-amazing-ledger/app/pagination"
)

type AccountEntryRequest struct {
	Account   AccountPath
	StartDate time.Time
	EndDate   time.Time
	Page      pagination.Page
}

type AccountEntryResponse struct {
	Entries  []AccountEntry
	NextPage pagination.Cursor
}

type AccountEntry struct {
	ID             uuid.UUID
	Version        Version
	Operation      OperationType
	Amount         int
	Event          int
	CompetenceDate time.Time
	Metadata       map[string]interface{}
}
