package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
)

func getAccountBalance(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting GetAccountBalance")
	defer log.Println("finishing GetAccountBalance")

	expectedBalance := 1000
	accountPathOne := "liability:stone:clients:" + uuid.New().String()
	accountPathTwo := "liability:stone:clients:" + uuid.New().String()

	// Define a new transaction with 2 entries
	t := conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, vos.NewAccountVersion, vos.CreditOperation, 1000)
	t.AddEntry(uuid.New(), accountPathTwo, vos.NewAccountVersion, vos.DebitOperation, 1000)
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
	t.AddEntry(uuid.New(), accountPathOne, vos.NewAccountVersion, vos.CreditOperation, 1000)
	t.AddEntry(uuid.New(), accountPathTwo, vos.NewAccountVersion, vos.DebitOperation, 1000)
	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	accountOne, err := conn.GetAccountBalance(context.Background(), accountPathOne)

	t = conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, accountOne.CurrentVersion(), vos.DebitOperation, 250)
	t.AddEntry(uuid.New(), accountPathTwo, vos.AnyAccountVersion, vos.CreditOperation, 250)
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
	AssertTrue(ledger.ErrAccountNotFound.Is(err))
}
