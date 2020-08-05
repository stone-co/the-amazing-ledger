package usecase

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/accounts/entities"
	"gotest.tools/assert"
)

func TestAccountsUseCase_Create(t *testing.T) {

	log := logrus.New()
	mockRepository := &accounts.RepositoryMock{}

	useCase := NewAccountUseCase(log, mockRepository)

	t.Run("Success Created", func(t *testing.T) {
		input := accounts.AccountInput{
			Type:  "asset",
			Owner: "str",
			Name:  "liquidacao",
		}

		mockRepository.OnCreate = func(account *entities.Account) error {
			return nil
		}

		err := useCase.CreateAccount(input)

		assert.Equal(t, nil, err)
	})

}
