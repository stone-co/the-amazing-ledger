package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func getAccountBalance(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting GetAccountBalance")
	defer log.Println("finishing GetAccountBalance")

	// Define a new transaction with 2 entries
	t := conn.NewTransaction(uuid.New())

	accountPathOne := "liability:stone:clients:" + uuid.New().String()
	accountPathTwo := "liability:stone:clients:" + uuid.New().String()
	t.AddEntry(uuid.New(), accountPathOne, entities.NewAccountVersion, entities.CreditOperation, 1000)
	t.AddEntry(uuid.New(), accountPathTwo, entities.NewAccountVersion, entities.DebitOperation, 1000)

	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	accountBalance, err := conn.GetAccountBalance(context.Background(), accountPathOne)

	AssertEqual(accountBalance.AccountName().Name(), accountPathOne)
	AssertEqual(accountBalance.Balance(), 1000)
	AssertEqual(nil, err)
}

func getAccountBalanceNotFoundAccount(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting getAccountBalanceNotFoundAccount")
	defer log.Println("finishing getAccountBalanceNotFoundAccount")

	accountPathNotFound := "liability:stone:clients:" + uuid.New().String()

	accountBalance, err := conn.GetAccountBalance(context.Background(), accountPathNotFound)

	AssertNil(accountBalance)
	AssertEqual(entities.ErrAccountNotFound, err)
}
