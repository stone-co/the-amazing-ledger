package transactions

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

type AccountInfoVersionAndBalance struct {
	AccountID      string `json:"account_id"`
	CurrentVersion uint64 `json:"current_version"`
	Balance        int    `json:"balance"`
}

type AccountNotFoundRequest struct {
	Message string `json:"message"`
}

func (h Handler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	log := h.log.WithFields(logrus.Fields{
		"handler": "Get Account Balance",
	})

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	accountID := vars["account_id"]
	decodedAccountID, err := url.QueryUnescape(accountID)
	if err != nil {
		log.WithError(err).Error("Query Unescape error")
		return
	}

	accountInfo, err := h.UseCase.GetAccountInfo(r.Context(), decodedAccountID)
	if err != nil {
		var messageJSON []byte

		if err == entities.ErrNotFound {
			messageJSON, _ = json.Marshal(AccountNotFoundRequest{Message: err.Error()})
			w.WriteHeader(http.StatusNotFound)
			log.WithError(err).Error("error getting account and balance")
		} else {
			messageJSON, _ = json.Marshal(AccountNotFoundRequest{Message: err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			log.WithError(err).Error("error getting account and balance")
		}

		_, err = w.Write(messageJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.WithError(err).Error("can't write response")
		}
		return
	}

	accountInfoVersionAndBalance := AccountInfoVersionAndBalance{
		AccountID:      accountInfo.AccountID,
		CurrentVersion: uint64(accountInfo.CurrentVersion.Current()),
		Balance:        accountInfo.Balance(),
	}

	accountInfoJSON, err := json.Marshal(accountInfoVersionAndBalance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(accountInfoJSON)
	if err != nil {
		log.WithError(err).Error("can't write response")
	}
}
