package entries

import (
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/entries"
)

type Handler struct {
	log     *logrus.Logger
	UseCase entries.UseCase
}

func NewEntriesHandler(log *logrus.Logger, useCase entries.UseCase) *Handler {
	return &Handler{
		log:     log,
		UseCase: useCase,
	}
}
