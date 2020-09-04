package usecase

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	"gotest.tools/assert"
)

func TestAccountsUseCase_Create(t *testing.T) {

	log := logrus.New()
	mockRepository := &ledger.RepositoryMock{}

	useCase := NewAccountsUseCase(log, mockRepository)

	t.Run("Successfully creates a sample account with minimum inputs", func(t *testing.T) {
		input := ledger.AccountInput{
			Type:  "asset",
			Owner: "owner",
			Name:  "account_name",
		}

		mockRepository.OnCreateAccount = func(account *entities.Account) error {
			return nil
		}

		err := useCase.CreateAccount(input)

		assert.Equal(t, nil, err)
	})

	t.Run("Successfully creates a sample account with a complete input", func(t *testing.T) {

		input := ledger.AccountInput{
			Type:     "asset",
			Owner:    "owner",
			Name:     "account_name",
			OwnerID:  uuid.New().String(),
			Metadata: []string{"metadatum_1", "metadatum_2"},
		}

		mockRepository.OnCreateAccount = func(account *entities.Account) error {
			return nil
		}

		err := useCase.CreateAccount(input)

		assert.Equal(t, nil, err)
	})

	t.Run("Fails to create an account when missing 'Type'", func(t *testing.T) {
		input := ledger.AccountInput{
			Owner: "owner",
			Name:  "account_name",
		}

		mockRepository.OnCreateAccount = func(account *entities.Account) error {
			return nil
		}

		err := useCase.CreateAccount(input)

		assert.Error(t, err, "missing 'type' input field")
	})

	t.Run("Fails to create an account when missing 'Owner'", func(t *testing.T) {
		input := ledger.AccountInput{
			Type: "asset",
			Name: "account_name",
		}

		mockRepository.OnCreateAccount = func(account *entities.Account) error {
			return nil
		}

		err := useCase.CreateAccount(input)

		assert.Error(t, err, "missing 'owner' input field")
	})

	t.Run("Fails to create an account when missing 'Name'", func(t *testing.T) {
		input := ledger.AccountInput{
			Type:  "asset",
			Owner: "owner",
		}

		mockRepository.OnCreateAccount = func(account *entities.Account) error {
			return nil
		}

		err := useCase.CreateAccount(input)

		assert.Error(t, err, "missing 'name' input field")
	})

	t.Run("Fails to creates an account with invalid 'Type'", func(t *testing.T) {
		input := ledger.AccountInput{
			Type:  "revenue",
			Owner: "owner",
			Name:  "account_name",
		}

		mockRepository.OnCreateAccount = func(account *entities.Account) error {
			return nil
		}

		err := useCase.CreateAccount(input)

		assert.Error(t, err, fmt.Sprintf("unknown account type '%s'", input.Type))
	})

}

func TestAccountsUseCase_Get(t *testing.T) {

	log := logrus.New()
	mockRepository := &ledger.RepositoryMock{}

	useCase := NewAccountsUseCase(log, mockRepository)

	t.Run("Successfully returns an account by id", func(t *testing.T) {
		id := uuid.New().String()
		now := time.Now()
		accountOutput := entities.Account{
			ID:        id,
			OwnerID:   "owner_id",
			Type:      entities.AccountType("asset"),
			Balance:   0,
			Owner:     "owner",
			Name:      "name",
			Metadata:  []string{"teste"},
			CreatedAt: now,
			UpdatedAt: nil,
		}

		mockRepository.OnGetAccount = func(id string) (entities.Account, error) {
			return accountOutput, nil
		}

		account, err := useCase.GetAccount(id)

		assert.Equal(t, nil, err)
		assert.Equal(t, accountOutput.ID, account.ID)
		assert.Equal(t, accountOutput.OwnerID, account.OwnerID)
		assert.Equal(t, string(accountOutput.Type), account.Type)
		assert.Equal(t, accountOutput.Balance, account.Balance)
		assert.Equal(t, accountOutput.Name, account.Name)
		assert.Equal(t, accountOutput.Owner, account.Owner)
		assert.Equal(t, accountOutput.Metadata[0], account.Metadata[0])
	})
}
