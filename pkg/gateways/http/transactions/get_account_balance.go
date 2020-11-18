package transactions

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
	"github.com/stone-co/the-amazing-ledger/pkg/gateways/http/responses"
)

type AccountBalanceResponse struct {
	AccountPath    string `json:"account_path"`
	CurrentVersion uint64 `json:"current_version"`
	TotalCredit    int    `json:"total_credit"`
	TotalDebit     int    `json:"total_debit"`
	Balance        int    `json:"balance"`
}

func (h Handler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	log := h.log.WithFields(logrus.Fields{
		"handler": "GetAccountBalance",
	})

	vars := mux.Vars(r)
	encodedAccountPath := vars["account_path"]
	decodedAccountPath, err := url.QueryUnescape(encodedAccountPath)
	if err != nil {
		log.WithError(err).Error("query unescape error")
		responses.SendError(w, log, err.Error(), http.StatusBadRequest)
		return
	}

	accountName, err := entities.NewAccountName(decodedAccountPath)
	if err != nil {
		log.WithError(err).Error("can't create account name")
		responses.SendError(w, log, err.Error(), http.StatusBadRequest)
		return
	}

	accountBalance, err := h.UseCase.GetAccountBalance(r.Context(), *accountName)
	if err != nil {
		if err == entities.ErrAccountNotFound {
			log.WithError(err).Error("account name does not exist")
			responses.SendError(w, log, err.Error(), http.StatusNotFound)
			return
		}

		log.WithError(err).Error("can't get account")
		responses.SendError(w, log, err.Error(), http.StatusBadRequest)
		return
	}

	accountBalanceResponse := AccountBalanceResponse{
		AccountPath:    accountBalance.AccountName.Name(),
		CurrentVersion: accountBalance.CurrentVersion.ToUInt64(),
		TotalCredit:    accountBalance.TotalCredit,
		TotalDebit:     accountBalance.TotalDebit,
		Balance:        accountBalance.Balance(),
	}

	responses.Send(w, log, accountBalanceResponse, http.StatusOK)
}
