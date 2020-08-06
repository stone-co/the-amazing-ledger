package usecase

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts/entities"
	"gotest.tools/assert"
)

func TestAccountsUseCase_Create(t *testing.T) {

	log := logrus.New()
	mockRepository := &accounts.RepositoryMock{}

	useCase := NewAccountUseCase(log, mockRepository)

	t.Run("Sucessfully creates a sample account with minimum inputs", func(t *testing.T) {
		input := accounts.AccountInput{
			Type:  "asset",
			Owner: "owner",
			Name:  "account_name",
		}

		mockRepository.OnCreate = func(account *entities.Account) error {
			return nil
		}

		err := useCase.CreateAccount(input)

		assert.Equal(t, nil, err)
	})

	t.Run("Sucessfully creates a sample account with a complete input", func(t *testing.T) {

		input := accounts.AccountInput{
			Type:     "asset",
			Owner:    "owner",
			Name:     "account_name",
			OwnerID:  uuid.New().String(),
			Metadata: []string{"metadatum_1", "metadatum_2"},
		}

		mockRepository.OnCreate = func(account *entities.Account) error {
			return nil
		}

		err := useCase.CreateAccount(input)

		assert.Equal(t, nil, err)
	})

	t.Run("Fails to create an account when missing 'Type'", func(t *testing.T) {
		input := accounts.AccountInput{
			Owner: "owner",
			Name:  "account_name",
		}

		mockRepository.OnCreate = func(account *entities.Account) error {
			return nil
		}

		err := useCase.CreateAccount(input)

		assert.Error(t, err, "missing 'type' input field")
	})

	t.Run("Fails to create an account when missing 'Owner'", func(t *testing.T) {
		input := accounts.AccountInput{
			Type: "asset",
			Name: "account_name",
		}

		mockRepository.OnCreate = func(account *entities.Account) error {
			return nil
		}

		err := useCase.CreateAccount(input)

		assert.Error(t, err, "missing 'owner' input field")
	})

	t.Run("Fails to create an account when missing 'Name'", func(t *testing.T) {
		input := accounts.AccountInput{
			Type:  "asset",
			Owner: "owner",
		}

		mockRepository.OnCreate = func(account *entities.Account) error {
			return nil
		}

		err := useCase.CreateAccount(input)

		assert.Error(t, err, "missing 'name' input field")
	})

	t.Run("Fails to creates an account with invalid 'Type'", func(t *testing.T) {
		input := accounts.AccountInput{
			Type:  "revenue",
			Owner: "owner",
			Name:  "account_name",
		}

		mockRepository.OnCreate = func(account *entities.Account) error {
			return nil
		}

		err := useCase.CreateAccount(input)

		assert.Error(t, err, fmt.Sprintf("unknown account type '%s'", input.Type))
	})

}
