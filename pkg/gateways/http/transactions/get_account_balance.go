package transactions

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

type AccountBalanceResponse struct {
	AccountPath    string `json:"account_path"`
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
	accountPath := vars["account_path"]
	decodedAccountPath, err := url.QueryUnescape(accountPath)
	if err != nil {
		log.WithError(err).Error("Query Unescape error")
		return
	}

	accountInfo, err := h.UseCase.GetAccountInfo(r.Context(), decodedAccountPath)
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

	accountBalanceResponse := AccountBalanceResponse{
		AccountPath:    accountInfo.AccountPath,
		CurrentVersion: uint64(accountInfo.CurrentVersion.Current()),
		Balance:        accountInfo.Balance(),
	}

	accountBalanceJSON, err := json.Marshal(accountBalanceResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(accountBalanceJSON)
	if err != nil {
		log.WithError(err).Error("can't write response")
	}
}
