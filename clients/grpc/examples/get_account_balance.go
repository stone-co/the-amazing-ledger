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

	expectedBalance := 1000
	accountPathOne := "liability:stone:clients:" + uuid.New().String()
	accountPathTwo := "liability:stone:clients:" + uuid.New().String()

	// Define a new transaction with 2 entries
	t := conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, entities.NewAccountVersion, entities.CreditOperation, 1000)
	t.AddEntry(uuid.New(), accountPathTwo, entities.NewAccountVersion, entities.DebitOperation, 1000)
	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	accountBalance, err := conn.GetAccountBalance(context.Background(), accountPathOne)

	AssertEqual(accountPathOne, accountBalance.AccountName().Name())
	AssertEqual(expectedBalance, accountBalance.Balance())
	AssertEqual(nil, err)
}

func getAccountBalanceWithMoreEntries(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting getAccountBalanceWithMoreEntries")
	defer log.Println("finishing getAccountBalanceWithMoreEntries")

	expectedBalance := 750
	accountPathOne := "liability:stone:clients:" + uuid.New().String()
	accountPathTwo := "liability:stone:clients:" + uuid.New().String()

	t := conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, entities.NewAccountVersion, entities.CreditOperation, 1000)
	t.AddEntry(uuid.New(), accountPathTwo, entities.NewAccountVersion, entities.DebitOperation, 1000)
	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	accountOne, err := conn.GetAccountBalance(context.Background(), accountPathOne)

	t = conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, accountOne.CurrentVersion(), entities.DebitOperation, 250)
	t.AddEntry(uuid.New(), accountPathTwo, entities.AnyAccountVersion, entities.CreditOperation, 250)
	err = conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	accountBalance, err := conn.GetAccountBalance(context.Background(), accountPathOne)

	AssertEqual(accountPathOne, accountBalance.AccountName().Name())
	AssertEqual(expectedBalance, accountBalance.Balance())
	AssertEqual(nil, err)
}

func getAccountBalanceNotFoundAccount(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting getAccountBalanceNotFoundAccount")
	defer log.Println("finishing getAccountBalanceNotFoundAccount")

	accountPathNotFound := "liability:stone:clients:" + uuid.New().String()

	accountBalance, err := conn.GetAccountBalance(context.Background(), accountPathNotFound)

	AssertNil(accountBalance)
	AssertTrue(entities.ErrAccountNotFound.Is(err))
}
