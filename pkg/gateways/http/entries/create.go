package entries

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/entries"
)

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	log := h.log.WithFields(logrus.Fields{
		"handler": "CreateTransaction",
	})

	var input []entries.EntryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.WithError(err).Error("can't decode request body into struct")
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.WithError(err).Error("can't write response")
		}
		return
	}

	if err := h.UseCase.CreateTransaction(input); err != nil {
		log.WithError(err).Error("error creating transaction")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.WithError(err).Error("can't write response")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}
