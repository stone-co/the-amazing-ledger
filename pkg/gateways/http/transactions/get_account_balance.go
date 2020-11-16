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
	TotalCredit    int    `json:"total_credit"`
	TotalDebit     int    `json:"total_debit"`
	Balance        int    `json:"balance"`
}

type AccountNotFoundRequest struct {
	Message string `json:"message"`
}

func (h Handler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	log := h.log.WithFields(logrus.Fields{
		"handler": "GetAccountBalance",
	})

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	encodedAccountPath := vars["account_path"]
	decodedAccountPath, err := url.QueryUnescape(encodedAccountPath)
	if err != nil {
		log.WithError(err).Error("Query Unescape error")
		return
	}

	accountName, err := entities.NewAccountName(decodedAccountPath)

	accountBalance, err := h.UseCase.GetAccountBalance(r.Context(), *accountName)
	if err != nil {
		var messageJSON []byte

		if err == entities.ErrAccountNotFound {
			messageJSON, _ = json.Marshal(AccountNotFoundRequest{Message: err.Error()})
			w.WriteHeader(http.StatusNotFound)
			log.WithError(err).Error("account name does not exist")
		} else {
			messageJSON, _ = json.Marshal(AccountNotFoundRequest{Message: err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			log.WithError(err).Error("can't get account")
		}

		_, err = w.Write(messageJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.WithError(err).Error("can't write http body response")
		}
		return
	}

	accountBalanceResponse := AccountBalanceResponse{
		AccountPath:    accountBalance.AccountName.Name(),
		CurrentVersion: uint64(accountBalance.CurrentVersion.Current()),
		TotalCredit:    accountBalance.TotalCredit,
		TotalDebit:     accountBalance.TotalDebit,
		Balance:        accountBalance.Balance(),
	}

	response, err := json.Marshal(accountBalanceResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.WithError(err).Error("can't write http body response")
	}
}
