package usecase

import (
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

func (t Transactions) CreateOperation(input transactions.OperationInput) error {
	return nil
}
