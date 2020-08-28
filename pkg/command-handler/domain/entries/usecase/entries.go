package usecase

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/entries"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/entries/entities"
)

type Entries struct {
	log        *logrus.Logger
	repository entries.Repository
}

func NewEntriesUseCase(log *logrus.Logger, repository entries.Repository) *Entries {
	return &Entries{
		log:        log,
		repository: repository,
	}
}

func (t Entries) CreateTransaction(input []entries.EntryInput) error {
	var err error = nil
	transaction := make([]entities.Entry, len(input))

	// check for empty or single element slice
	if len(input) == 0 {
		return errors.New("empty input")
	}
	if len(input) == 1 {
		return errors.New("transaction needs at least 2 entries")
	}

	var sumAmount int = 0
	// check for required fields and build entity
	for i, t := range input {
		if t.AccountType == "" {
			err = errors.New("missing 'account_type' input field")
			break
		}
		if t.AccountOwner == "" {
			err = errors.New("missing 'account_owner' input field")
			break
		}
		if t.AccountName == "" {
			err = errors.New("missing 'account_name' input field")
			break
		}
		if t.AccountMetadata == nil {
			t.AccountMetadata = []string{}
		}
		if t.Amount == 0 {
			err = errors.New("amount cannot be 0")
			break
		}
		sumAmount += t.Amount
		entry := entities.Entry{
			RequestID: t.RequestID,
			Amount:    t.Amount,
		}
		transaction[i] = entry
	}

	// check for zero-sum
	if sumAmount != 0 {
		err = errors.New("sum of all amounts must be 0!")
	}

	//TO-DO check for valid accounts

	// insert transaction atomically in database
	if create_err := t.repository.Create(&transaction); create_err != nil {
		return fmt.Errorf("can't create entries: %s", create_err.Error())
	}

	return err
}
