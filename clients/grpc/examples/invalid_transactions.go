package main

import (
	"log"

	"github.com/google/uuid"

	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
	"github.com/stone-co/the-amazing-ledger/pkg/command-handler/domain/ledger/entities"
)

func invalidTransactionsTests(conn *ledger.Connection) {
	transactionWithInvalidIdReturnsInvalidData(conn)
	entryWithInvalidIdReturnsInvalidData(conn)
	withoutEntriesReturnsInvalidEntriesNumber(conn)
}

func transactionWithInvalidIdReturnsInvalidData(conn *ledger.Connection) {
	log.Println("transactionWithInvalidIdReturnsInvalidData example starting...")
	defer log.Println("transactionWithInvalidIdReturnsInvalidData example finishing...")

	invalidUUID := uuid.Nil
	t := conn.NewTransaction(invalidUUID)

	accountID1 := uuid.New().String()
	accountID2 := uuid.New().String()

	t.AddEntry(uuid.New(), accountID1, entities.NewAccountVersion, ledger.Debit, 15000)
	t.AddEntry(uuid.New(), accountID2, entities.NewAccountVersion, ledger.Credit, 15000)

	err := conn.SaveTransaction(t)
	AssertEqual(entities.ErrInvalidData, err)
}

func entryWithInvalidIdReturnsInvalidData(conn *ledger.Connection) {
	log.Println("entryWithInvalidIdReturnsInvalidData example starting...")
	defer log.Println("entryWithInvalidIdReturnsInvalidData example finishing...")

	t := conn.NewTransaction(uuid.New())

	accountID1 := uuid.New().String()
	accountID2 := uuid.New().String()

	invalidUUID := uuid.Nil
	t.AddEntry(uuid.New(), accountID1, entities.NewAccountVersion, ledger.Debit, 15000)
	t.AddEntry(invalidUUID, accountID2, entities.NewAccountVersion, ledger.Credit, 15000)

	err := conn.SaveTransaction(t)
	AssertEqual(entities.ErrInvalidData, err)
}

func withoutEntriesReturnsInvalidEntriesNumber(conn *ledger.Connection) {
	log.Println("entryWithInvalidIdReturnsInvalidData example starting...")
	defer log.Println("entryWithInvalidIdReturnsInvalidData example finishing...")

	t := conn.NewTransaction(uuid.New())

	err := conn.SaveTransaction(t)
	AssertEqual(entities.ErrInvalidEntriesNumber, err)
}
