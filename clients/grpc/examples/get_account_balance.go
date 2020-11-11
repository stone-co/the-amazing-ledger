package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

var accountPathOne = "liability:stone:clients:" + uuid.New().String()
var accountPathNotFound = "liability:stone:clients:" + uuid.New().String()

func saveTransactionToGetAccountTest(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting saveTransactionToGetAccountTest")
	defer log.Println("finishing saveTransactionToGetAccountTest")

	// Define a new transaction with 3 entries
	t := conn.NewTransaction(uuid.New())

	accountPathTwo := "liability:stone:clients:" + uuid.New().String()
	t.AddEntry(uuid.New(), accountPathOne, entities.NewAccountVersion, entities.CreditOperation, 1000)
	t.AddEntry(uuid.New(), accountPathTwo, entities.NewAccountVersion, entities.DebitOperation, 1000)

	err := conn.SaveTransaction(context.Background(), t)

	AssertEqual(nil, err)
}

func getAccountBalance(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting GetAccountBalance")
	defer log.Println("finishing GetAccountBalance")

	saveTransactionToGetAccountTest(log, conn)

	a := conn.NewAccountRequest(accountPathOne)

	accountBalance, err := conn.GetAccountBalance(context.Background(), a)

	AssertEqual(accountBalance.Balance, 1000)
	AssertEqual(nil, err)
}

func getAccountBalanceNotFoundAccount(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting getAccountBalanceNotFoundAccount")
	defer log.Println("finishing getAccountBalanceNotFoundAccount")

	a := conn.NewAccountRequest(accountPathNotFound)

	_, err := conn.GetAccountBalance(context.Background(), a)

	AssertEqual(entities.ErrNotFound, err)
}
