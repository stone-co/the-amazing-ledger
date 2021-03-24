package main

import (
	"context"
	"strings"

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

func getAccountBalanceWithWildcard(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting getAccountBalanceWithWildcard")
	defer log.Println("finishing getAccountBalanceWithWildcard")

	transactionOne := 1100
	transactionTwo := 900
	transactionThree := 1000
	expectedBalance := transactionOne + transactionTwo + transactionThree

	newUUID := uuid.New().String()
	accountPathOne := "liability:stone:clients:" + newUUID + "/test_1"
	accountPathTwo := "liability:stone:clients:" + newUUID + "/test_2"
	accountPathThree := "liability:stone:clients:" + newUUID
	accountPathFour := "liability:stone:clients:" + uuid.New().String()

	t := conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, vos.NewAccountVersion, vos.CreditOperation, transactionOne)
	t.AddEntry(uuid.New(), accountPathFour, vos.NewAccountVersion, vos.DebitOperation, transactionOne)

	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	accountFour, err := conn.GetAccountBalance(context.Background(), accountPathFour)
	AssertEqual(nil, err)

	t = conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathTwo, vos.NewAccountVersion, vos.CreditOperation, transactionTwo)
	t.AddEntry(uuid.New(), accountPathFour, accountFour.CurrentVersion(), vos.DebitOperation, transactionTwo)

	err = conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	accountFour, err = conn.GetAccountBalance(context.Background(), accountPathFour)
	AssertEqual(nil, err)

	t = conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathThree, vos.NewAccountVersion, vos.CreditOperation, transactionThree)
	t.AddEntry(uuid.New(), accountPathFour, accountFour.CurrentVersion(), vos.DebitOperation, transactionThree)

	err = conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	// Check wildcard balance is ok
	wildcardPath := strings.Replace(accountPathOne, "/test_1", "/*", 1)
	wildcard, err := conn.GetAccountBalance(context.Background(), wildcardPath)
	AssertEqual(nil, err)

	AssertEqual(wildcardPath, wildcard.AccountName().Name())
	AssertEqual(expectedBalance, wildcard.Balance())
	AssertEqual(nil, err)

	// Check separate account balances are ok
	accountOne, err := conn.GetAccountBalance(context.Background(), accountPathOne)
	AssertEqual(nil, err)

	AssertEqual(accountPathOne, accountOne.AccountName().Name())
	AssertEqual(transactionOne, accountOne.Balance())
	AssertEqual(nil, err)

	accountTwo, err := conn.GetAccountBalance(context.Background(), accountPathTwo)
	AssertEqual(nil, err)

	AssertEqual(accountPathTwo, accountTwo.AccountName().Name())
	AssertEqual(transactionTwo, accountTwo.Balance())
	AssertEqual(nil, err)

	accountThree, err := conn.GetAccountBalance(context.Background(), accountPathThree)
	AssertEqual(nil, err)

	AssertEqual(accountPathThree, accountThree.AccountName().Name())
	AssertEqual(transactionThree, accountThree.Balance())
	AssertEqual(nil, err)
}
