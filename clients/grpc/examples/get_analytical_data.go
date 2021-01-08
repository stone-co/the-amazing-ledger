package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
)

// Scenarios with a valid path.
func getAnalyticalData(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting getAnalyticalData")
	defer log.Println("finishing getAnalyticalData")

	// Without entries
	getAnalyticalDataWithoutEntries(log, conn)

	// With 1 entry
	getAnalyticalDataWithOneEntry(log, conn)

	// With 3 entries for the same path
	getAnalyticalDataWithThreeEntries(log, conn)

	// With a lot of entries, and some match (and some no) with the path
	getAnalyticalDataWithPartialPath(log, conn)
}

func getAnalyticalDataWithoutEntries(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("running getAnalyticalDataWithoutEntries")

	entries, err := conn.GetAnalyticalData(context.Background(), "liability:stone:"+uuid.New().String())
	AssertEqual(nil, err)
	AssertTrue(entries != nil)
	AssertTrue(len(entries) == 0)
}

func getAnalyticalDataWithOneEntry(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("running getAnalyticalDataWithOneEntry")

	path := "liability:stone:" + uuid.New().String()
	account1 := "liability:stone:clients:" + uuid.New().String()
	account2 := path + ":" + uuid.New().String()

	// Define a new transaction.
	t := conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), account1, vos.NewAccountVersion, vos.CreditOperation, 1000)
	t.AddEntry(uuid.New(), account2, vos.NewAccountVersion, vos.DebitOperation, 1000)
	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	entries, err := conn.GetAnalyticalData(context.Background(), path)
	AssertEqual(nil, err)
	AssertTrue(len(entries) == 1)

	entry := entries[0]
	AssertEqual(account2, entry.Account)
	AssertEqual(vos.DebitOperation, entry.Operation)
	AssertEqual(1000, entry.Amount)
}

func getAnalyticalDataWithThreeEntries(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("running getAnalyticalDataWithThreeEntries")

	account1 := "liability:stone:clients:" + uuid.New().String()

	path := "liability:stone:" + uuid.New().String()
	account2 := path + ":" + uuid.New().String()
	account3 := path + ":" + uuid.New().String()

	// Create a new transaction.
	t := conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), account1, vos.NewAccountVersion, vos.CreditOperation, 1000)
	t.AddEntry(uuid.New(), account2, vos.NewAccountVersion, vos.DebitOperation, 1000)
	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	// Create a new transaction.
	t = conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), account1, vos.AnyAccountVersion, vos.DebitOperation, 2000)
	t.AddEntry(uuid.New(), account2, vos.AnyAccountVersion, vos.CreditOperation, 1500)
	t.AddEntry(uuid.New(), account3, vos.AnyAccountVersion, vos.CreditOperation, 500)
	err = conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	entries, err := conn.GetAnalyticalData(context.Background(), path)
	AssertEqual(nil, err)
	AssertTrue(len(entries) == 3)

	entry := entries[0]
	AssertEqual(account2, entry.Account)
	AssertEqual(vos.DebitOperation, entry.Operation)
	AssertEqual(1000, entry.Amount)

	entry = entries[1]
	AssertEqual(account2, entry.Account)
	AssertEqual(vos.CreditOperation, entry.Operation)
	AssertEqual(1500, entry.Amount)

	entry = entries[2]
	AssertEqual(account3, entry.Account)
	AssertEqual(vos.CreditOperation, entry.Operation)
	AssertEqual(500, entry.Amount)
}

func getAnalyticalDataWithPartialPath(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("running getAnalyticalDataWithPartialPath")

	path := "liability:" + uuid.New().String()
	account1 := "liability:stone:clients:" + uuid.New().String()
	account2 := path + ":" + "aaa" + ":" + uuid.New().String()
	account3 := path + ":" + "bbb" + ":" + uuid.New().String()

	// Create a new transaction.
	t := conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), account1, vos.AnyAccountVersion, vos.CreditOperation, 2000)
	t.AddEntry(uuid.New(), account2, vos.AnyAccountVersion, vos.DebitOperation, 1500)
	t.AddEntry(uuid.New(), account3, vos.AnyAccountVersion, vos.DebitOperation, 500)
	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	entries, err := conn.GetAnalyticalData(context.Background(), path)
	AssertEqual(nil, err)
	AssertTrue(len(entries) == 2)

	entry := entries[0]
	AssertEqual(account2, entry.Account)
	AssertEqual(vos.DebitOperation, entry.Operation)
	AssertEqual(1500, entry.Amount)

	entry = entries[1]
	AssertEqual(account3, entry.Account)
	AssertEqual(vos.DebitOperation, entry.Operation)
	AssertEqual(500, entry.Amount)
}
