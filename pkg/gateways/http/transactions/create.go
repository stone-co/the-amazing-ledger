package transactions

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

type CreateTransactionRequest struct {
	ID      uuid.UUID `json:"id"`
	Entries []struct {
		ID        uuid.UUID `json:"id"`
		Operation string    `json:"operation"`
		AccountID uuid.UUID `json:"account_id"`
		Version   uint64    `json:"version"`
		Amount    int       `json:"amount"`
	} `json:"entries"`
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	log := h.log.WithFields(logrus.Fields{
		"handler": "CreateTransaction",
	})

	var input CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.WithError(err).Error("can't decode request body into struct")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.WithError(err).Error("can't write response")
		}
		return
	}

	entries := []entities.Entry{}
	for _, entry := range input.Entries {
		entries = append(entries, entities.Entry{
			ID:        entry.ID,
			Operation: entities.OperationTypeFromString(entry.Operation),
			AccountID: entry.AccountID,
			Version:   entities.VersionType(entry.Version),
			Amount:    entry.Amount,
		})
	}

	if err := h.UseCase.CreateTransaction(r.Context(), input.ID, entries); err != nil {
		log.WithError(err).Error("error creating transaction")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.WithError(err).Error("can't write response")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}
