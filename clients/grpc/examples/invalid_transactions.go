package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func invalidTransactionsTests(log *logrus.Entry, conn *ledger.Connection) {
	transactionWithInvalidIdReturnsInvalidData(log, conn)
	entryWithInvalidIdReturnsInvalidData(log, conn)
	withoutEntriesReturnsInvalidEntriesNumber(log, conn)
}

func transactionWithInvalidIdReturnsInvalidData(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting transactionWithInvalidIdReturnsInvalidData")
	defer log.Println("finishing transactionWithInvalidIdReturnsInvalidData")

	invalidUUID := uuid.Nil
	t := conn.NewTransaction(invalidUUID)

	accountID1 := "liability:clients:available:" + uuid.New().String()
	accountID2 := "liability:clients:available:" + uuid.New().String()

	t.AddEntry(uuid.New(), accountID1, entities.NewAccountVersion, entities.DebitOperation, 15000)
	t.AddEntry(uuid.New(), accountID2, entities.NewAccountVersion, entities.CreditOperation, 15000)

	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(entities.ErrInvalidData, err)
}

func entryWithInvalidIdReturnsInvalidData(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting entryWithInvalidIdReturnsInvalidData")
	defer log.Println("finishing entryWithInvalidIdReturnsInvalidData")

	t := conn.NewTransaction(uuid.New())

	accountID1 := "liability:clients:available:" + uuid.New().String()
	accountID2 := "liability:clients:available:" + uuid.New().String()

	invalidUUID := uuid.Nil
	t.AddEntry(uuid.New(), accountID1, entities.NewAccountVersion, entities.DebitOperation, 15000)
	t.AddEntry(invalidUUID, accountID2, entities.NewAccountVersion, entities.CreditOperation, 15000)

	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(entities.ErrInvalidData, err)
}

func withoutEntriesReturnsInvalidEntriesNumber(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting withoutEntriesReturnsInvalidEntriesNumber")
	defer log.Println("finishing withoutEntriesReturnsInvalidEntriesNumber")

	t := conn.NewTransaction(uuid.New())

	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(entities.ErrInvalidEntriesNumber, err)
}
