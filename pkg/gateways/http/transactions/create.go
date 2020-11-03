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
		ID              uuid.UUID `json:"id"`
		Operation       string    `json:"operation"`
		AccountID       string    `json:"account_id"`
		ExpectedVersion uint64    `json:"expected_version"`
		Amount          int       `json:"amount"`
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

	domainEntries := []entities.Entry{}
	for _, entry := range input.Entries {
		domainEntry, err := entities.NewEntry(
			entry.ID,
			entities.OperationTypeFromString(entry.Operation),
			entry.AccountID,
			entities.Version(entry.ExpectedVersion),
			entry.Amount,
		)
		if err != nil {
			log.WithError(err).Error("invalid entry data")
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.WithError(err).Error("can't write response")
			}
			return
		}
		domainEntries = append(domainEntries, *domainEntry)
	}

	if err := h.UseCase.CreateTransaction(r.Context(), input.ID, domainEntries); err != nil {
		log.WithError(err).Error("creating transaction")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.WithError(err).Error("can't write response")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}
