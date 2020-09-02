package entries

import (
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
)

type Handler struct {
	log     *logrus.Logger
	UseCase ledger.EntriesUseCase
}

func NewEntriesHandler(log *logrus.Logger, useCase ledger.EntriesUseCase) *Handler {
	return &Handler{
		log:     log,
		UseCase: useCase,
	}
}
