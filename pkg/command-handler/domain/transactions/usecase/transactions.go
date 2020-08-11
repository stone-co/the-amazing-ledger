package usecase

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/transactions"
)

type Transactions struct {
	log        *logrus.Logger
	repository transactions.Repository
}

func NewTransactionUseCase(log *logrus.Logger, repository transactions.Repository) *Transactions {
	return &Transactions{
		log:        log,
		repository: repository,
	}
}

func (t Transactions) CreateOperation(input []transactions.TransactionInput) error {
	var err error = nil

	// check for empty or single element slice
	if len(input) == 0 {
		return errors.New("empty input")
	}
	if len(input) == 1 {
		return errors.New("operation needs at least 2 transactions")
	}

	var sumAmount int = 0
	// check for required fields
	for _, t := range input {
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
	}

	// check for zero-sum
	if sumAmount != 0 {
		err = errors.New("sum of all amounts must be 0!")
	}

	// if err := t.repository.Create()

	return err
}
