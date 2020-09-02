package accounts

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger"
)

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	log := h.log.WithFields(logrus.Fields{
		"handler": "CreateAccount",
	})

	var input ledger.AccountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.WithError(err).Error("can't decode request body into struct")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.UseCase.CreateAccount(input); err != nil {
		log.WithError(err).Error("error creating account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
