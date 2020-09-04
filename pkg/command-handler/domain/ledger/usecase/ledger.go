package usecase

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

type LedgerUseCase struct {
	log        *logrus.Logger
	repository ledger.Repository
}

func NewLedgerUseCase(log *logrus.Logger, repository ledger.Repository) *LedgerUseCase {
	return &LedgerUseCase{
		log:        log,
		repository: repository,
	}
}

func newAccount(input ledger.AccountInput) (entities.Account, error) {
	account := entities.Account{}

	if input.Type == "" {
		return account, errors.New("missing 'type' input field")
	}
	if input.Owner == "" {
		return account, errors.New("missing 'owner' input field")
	}
	if input.Name == "" {
		return account, errors.New("missing 'name' input field")
	}
	if input.Metadata == nil {
		input.Metadata = []string{}
	}

	accountType := entities.AccountType(input.Type)
	if accountType != entities.Asset && accountType != entities.Liability {
		return account, fmt.Errorf("unknown account type '%s'", input.Type)
	}

	account = entities.Account{
		Type:     accountType,
		Owner:    input.Owner,
		Name:     input.Name,
		OwnerID:  input.OwnerID,
		Metadata: input.Metadata,
		Balance:  0,
	}

	return account, nil
}

func (l LedgerUseCase) CreateAccount(input ledger.AccountInput) (ledger.Account, error) {
	account, err := newAccount(input)

	if err != nil {
		return ledger.Account{}, err
	}

	return l.doCreateAccount(account)
}

func (l LedgerUseCase) doCreateAccount(account entities.Account) (ledger.Account, error) {
	account, err := l.repository.CreateAccount(&account)
	if err != nil {
		return ledger.Account{}, fmt.Errorf("can't create account: %s", err.Error())
	}

	newAccount := ledger.Account{
		ID:       account.ID,
		OwnerID:  account.OwnerID,
		Type:     string(account.Type),
		Balance:  account.Balance,
		Owner:    account.Owner,
		Name:     account.Name,
		Metadata: account.Metadata,
	}
	return newAccount, err
}

func (l LedgerUseCase) GetAccount(id string) (ledger.Account, error) {
	account, err := l.repository.GetAccount(id)
	if err != nil {
		var account = ledger.Account{}
		return account, fmt.Errorf("Can't get account of id %s: %s", id, err.Error())
	}

	newAccount := ledger.Account{
		ID:       account.ID,
		OwnerID:  account.OwnerID,
		Type:     string(account.Type),
		Balance:  account.Balance,
		Owner:    account.Owner,
		Name:     account.Name,
		Metadata: account.Metadata,
	}
	return newAccount, err
}

func (l LedgerUseCase) SearchAccount(input ledger.AccountInput) (ledger.Account, error) {
	account, err := newAccount(input)

	if err != nil {
		return ledger.Account{}, err
	}

	account, err = l.repository.SearchAccount(&account)
	if err != nil {
		return ledger.Account{}, fmt.Errorf("Can't get account: %s", err.Error())
	}

	newAccount := ledger.Account{
		ID:       account.ID,
		OwnerID:  account.OwnerID,
		Type:     string(account.Type),
		Balance:  account.Balance,
		Owner:    account.Owner,
		Name:     account.Name,
		Metadata: account.Metadata,
	}
	return newAccount, err
}

func (l LedgerUseCase) SearchOrCreateAccount(input ledger.AccountInput) (ledger.Account, error) {
	acc, err := l.SearchAccount(input)

	if err == nil {
		return acc, err
	}

	if acc.ID == "" {
		if acc, err = l.CreateAccount(input); err != nil {
			return acc, err
		}
	}

	return acc, err
}

func (l LedgerUseCase) UpdateBalance(id string, balance int) error {
	if err := l.repository.UpdateBalance(id, balance); err != nil {
		return fmt.Errorf("can't update balance: %s", err.Error())
	}

	return nil
}

func (l LedgerUseCase) CreateTransaction(input []ledger.EntryInput) error {
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
		account, err := l.SearchOrCreateAccount(ledger.AccountInput{
			Type:     e.AccountType,
			OwnerID:  e.AccountOwnerID,
			Owner:    e.AccountOwner,
			Name:     e.AccountName,
			Metadata: e.AccountMetadata,
		})

		if err != nil {
			return err
		}

		if e.Amount == 0 {
			return errors.New("amount cannot be 0")
		}
		sumAmount += e.Amount
		entry := entities.Entry{
			AccountID: account.ID,
			RequestID: e.RequestID,
			Amount:    e.Amount,
		}
		transaction[i] = entry
	}

	// check for zero-sum
	if sumAmount != 0 {
		return errors.New("sum of all amounts must be 0!")
	}

	// insert transaction atomically in database
	if create_err := l.repository.CreateTransaction(&transaction); create_err != nil {
		return fmt.Errorf("can't create entries: %s", create_err.Error())
	}

	return nil
}
