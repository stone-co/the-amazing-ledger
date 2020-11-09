package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

var accountID1 = "liability:stone:clients:" + uuid.New().String()
var accountIDNotFound = "liability:stone:clients:" + uuid.New().String()

func saveTransactionToGetAccountTest(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting saveTransactionToGetAccountTest")
	defer log.Println("finishing saveTransactionToGetAccountTest")

	// Define a new transaction with 3 entries
	t := conn.NewTransaction(uuid.New())

	accountID2 := "liability:stone:clients:" + uuid.New().String()
	t.AddEntry(uuid.New(), accountID1, entities.NewAccountVersion, entities.CreditOperation, 1000)
	t.AddEntry(uuid.New(), accountID2, entities.NewAccountVersion, entities.DebitOperation, 1000)

	err := conn.SaveTransaction(context.Background(), t)

	AssertEqual(nil, err)
}

func getAccountBalance(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting GetAccountBalance")
	defer log.Println("finishing GetAccountBalance")

	saveTransactionToGetAccountTest(log, conn)

	a := conn.NewAccountID(accountID1)

	accountBalance, err := conn.GetAccountBalance(context.Background(), a)

	AssertEqual(accountBalance.Balance, 1000)
	AssertEqual(nil, err)
}

func getAccountBalanceNotFoundAccount(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting getAccountBalanceNotFoundAccount")
	defer log.Println("finishing getAccountBalanceNotFoundAccount")

	a := conn.NewAccountID(accountIDNotFound)

	_, err := conn.GetAccountBalance(context.Background(), a)

	AssertEqual(entities.ErrNotFound, err)
}
