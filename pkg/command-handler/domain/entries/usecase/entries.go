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

func (e Entries) CreateTransaction(input []entries.EntryInput) error {
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
	for i, e := range input {
		if e.AccountType == "" {
			err = errors.New("missing 'account_type' input field")
			return err
		}
		if e.AccountOwner == "" {
			err = errors.New("missing 'account_owner' input field")
			return err
		}
		if e.AccountName == "" {
			err = errors.New("missing 'account_name' input field")
			return err
		}
		if e.AccountMetadata == nil {
			e.AccountMetadata = []string{}
		}
		if e.Amount == 0 {
			err = errors.New("amount cannot be 0")
			return err
		}
		sumAmount += e.Amount
		entry := entities.Entry{
			RequestID: e.RequestID,
			Amount:    e.Amount,
		}
		transaction[i] = entry
	}

	// check for zero-sum
	if sumAmount != 0 {
		return errors.New("sum of all amounts must be 0!")
	}

	//TO-DO check for valid accounts

	// insert transaction atomically in database
	if create_err := e.repository.Create(&transaction); create_err != nil {
		return fmt.Errorf("can't create entries: %s", create_err.Error())
	}

	return err
}
