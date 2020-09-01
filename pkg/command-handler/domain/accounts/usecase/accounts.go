package usecase

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts/entities"
)

type Accounts struct {
	log        *logrus.Logger
	repository accounts.Repository
}

func NewAccountsUseCase(log *logrus.Logger, repository accounts.Repository) *Accounts {
	return &Accounts{
		log:        log,
		repository: repository,
	}
}

func (a Accounts) CreateAccount(input accounts.AccountInput) error {
	if input.Type == "" {
		return errors.New("missing 'type' input field")
	}
	if input.Owner == "" {
		return errors.New("missing 'owner' input field")
	}
	if input.Name == "" {
		return errors.New("missing 'name' input field")
	}
	if input.Metadata == nil {
		input.Metadata = []string{}
	}

	accountType := entities.AccountType(input.Type)
	if accountType != entities.Asset && accountType != entities.Liability {
		return fmt.Errorf("unknown account type '%s'", input.Type)
	}

	account := entities.Account{
		Type:     accountType,
		Owner:    input.Owner,
		Name:     input.Name,
		OwnerID:  input.OwnerID,
		Metadata: input.Metadata,
		Balance:  0,
	}

	if err := a.repository.Create(&account); err != nil {
		return fmt.Errorf("can't create account: %s", err.Error())
	}

	return nil
}

func (a Accounts) GetAccount(id string) (accounts.Account, error) {
	account, err := a.repository.Get(id)
	if err != nil {
		var account = accounts.Account{}
		return account, fmt.Errorf("Can't get account of id %s: %s", id, err.Error())
	}

	newAccount := accounts.Account{
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

func (a Accounts) SearchAccount(accountType string, accountOwnerID string, accountOwner string, accountName string, accountMetadata string) (accounts.Account, error) {
	account, err := a.repository.Search(accountType, accountOwnerID, accountOwner, accountName, accountMetadata)
	if err != nil {
		var account = accounts.Account{}
		return account, fmt.Errorf("Can't get account: %s", err.Error())
	}

	newAccount := accounts.Account{
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

func (a Accounts) UpdateBalance(id string, balance int) error {
	if err := a.repository.Update(id, balance); err != nil {
		return fmt.Errorf("can't update balance: %s", err.Error())
	}

	return nil
}
