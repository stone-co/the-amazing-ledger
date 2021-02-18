package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
)

func getAccountHistory(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting GetAccountHistory")
	defer log.Println("finishing GetAccountHistory")

	accountPathOne := "liability:stone:clients:" + uuid.New().String()
	accountPathTwo := "liability:stone:clients:" + uuid.New().String()

	t := conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, vos.NewAccountVersion, vos.CreditOperation, 1000)
	t.AddEntry(uuid.New(), accountPathTwo, vos.NewAccountVersion, vos.DebitOperation, 1000)
	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	accountOne, err := conn.GetAccountBalance(context.Background(), accountPathOne)

	t = conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, accountOne.CurrentVersion(), vos.DebitOperation, 500)
	t.AddEntry(uuid.New(), accountPathTwo, vos.AnyAccountVersion, vos.CreditOperation, 500)
	err = conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	accountHistoryOne, err := conn.GetAccountHistory(context.Background(), accountPathOne)
	AssertEqual(nil, err)

	AssertEqual(2, len(accountHistoryOne))

	AssertEqual(1000, accountHistoryOne[0].Amount())
	AssertEqual(vos.CreditOperation, accountHistoryOne[0].Operation())

	AssertEqual(500, accountHistoryOne[1].Amount())
	AssertEqual(vos.DebitOperation, accountHistoryOne[1].Operation())

	accountHistoryTwo, err := conn.GetAccountHistory(context.Background(), accountPathTwo)
	AssertEqual(nil, err)

	AssertEqual(2, len(accountHistoryTwo))

	AssertEqual(1000, accountHistoryTwo[0].Amount())
	AssertEqual(vos.DebitOperation, accountHistoryTwo[0].Operation())

	AssertEqual(500, accountHistoryTwo[1].Amount())
	AssertEqual(vos.CreditOperation, accountHistoryTwo[1].Operation())
}

func getAccountHistoryWithForEntries(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting getAccountHistoryWithForEntries")
	defer log.Println("finishing getAccountHistoryWithForEntries")

	accountPathOne := "liability:stone:clients:" + uuid.New().String()
	accountPathTwo := "liability:stone:clients:" + uuid.New().String()

	t := conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, vos.NewAccountVersion, vos.CreditOperation, 1000)
	t.AddEntry(uuid.New(), accountPathTwo, vos.NewAccountVersion, vos.DebitOperation, 1000)
	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	accountOne, err := conn.GetAccountBalance(context.Background(), accountPathOne)

	t = conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, accountOne.CurrentVersion(), vos.CreditOperation, 500)
	t.AddEntry(uuid.New(), accountPathTwo, vos.AnyAccountVersion, vos.DebitOperation, 500)
	err = conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	accountOne, err = conn.GetAccountBalance(context.Background(), accountPathOne)

	t = conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, accountOne.CurrentVersion(), vos.DebitOperation, 500)
	t.AddEntry(uuid.New(), accountPathTwo, vos.AnyAccountVersion, vos.CreditOperation, 500)
	err = conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	accountOne, err = conn.GetAccountBalance(context.Background(), accountPathOne)

	t = conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, accountOne.CurrentVersion(), vos.DebitOperation, 1000)
	t.AddEntry(uuid.New(), accountPathTwo, vos.AnyAccountVersion, vos.CreditOperation, 1000)
	err = conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	accountHistoryOne, err := conn.GetAccountHistory(context.Background(), accountPathOne)
	AssertEqual(nil, err)

	AssertEqual(4, len(accountHistoryOne))

	AssertEqual(1000, accountHistoryOne[0].Amount())
	AssertEqual(vos.CreditOperation, accountHistoryOne[0].Operation())

	AssertEqual(500, accountHistoryOne[1].Amount())
	AssertEqual(vos.CreditOperation, accountHistoryOne[1].Operation())

	AssertEqual(500, accountHistoryOne[2].Amount())
	AssertEqual(vos.DebitOperation, accountHistoryOne[2].Operation())

	AssertEqual(1000, accountHistoryOne[3].Amount())
	AssertEqual(vos.DebitOperation, accountHistoryOne[3].Operation())

	accountHistoryTwo, err := conn.GetAccountHistory(context.Background(), accountPathTwo)
	AssertEqual(nil, err)

	AssertEqual(4, len(accountHistoryTwo))

	AssertEqual(1000, accountHistoryTwo[0].Amount())
	AssertEqual(vos.DebitOperation, accountHistoryTwo[0].Operation())

	AssertEqual(500, accountHistoryTwo[1].Amount())
	AssertEqual(vos.DebitOperation, accountHistoryTwo[1].Operation())

	AssertEqual(500, accountHistoryTwo[2].Amount())
	AssertEqual(vos.CreditOperation, accountHistoryTwo[2].Operation())

	AssertEqual(1000, accountHistoryTwo[3].Amount())
	AssertEqual(vos.CreditOperation, accountHistoryTwo[3].Operation())
}
