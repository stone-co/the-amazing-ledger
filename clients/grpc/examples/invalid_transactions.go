package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
)

func invalidTransactionsTests(log *logrus.Entry, conn *ledger.Connection) {
	transactionWithInvalidIdReturnsInvalidTransactionID(log, conn)
	entryWithInvalidIdReturnsInvalidEntryID(log, conn)
	withoutEntriesReturnsInvalidEntriesNumber(log, conn)
	entryWithInvalidVersion(log, conn)
	entryWithInvalidAccountStructure(log, conn)
	transactionWithInvalidBalanceReturnsErrInvalidBalance(log, conn)
	transactionWithInvalidIdempotencyKeyReturnsErrIdempotencyKeyViolation(log, conn)
}

func transactionWithInvalidIdReturnsInvalidTransactionID(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting transactionWithInvalidIdReturnsInvalidTransactionID")
	defer log.Println("finishing transactionWithInvalidIdReturnsInvalidTransactionID")

	invalidUUID := uuid.Nil
	t := conn.NewTransaction(invalidUUID)

	accountID1 := "liability:clients:available:" + uuid.New().String()
	accountID2 := "liability:clients:available:" + uuid.New().String()

	t.AddEntry(uuid.New(), accountID1, vos.NewAccountVersion, vos.DebitOperation, 15000)
	t.AddEntry(uuid.New(), accountID2, vos.NewAccountVersion, vos.CreditOperation, 15000)

	err := conn.SaveTransaction(context.Background(), t)
	AssertTrue(ledger.ErrInvalidTransactionID.Is(err))
}

func entryWithInvalidIdReturnsInvalidEntryID(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting entryWithInvalidIdReturnsInvalidEntryID")
	defer log.Println("finishing entryWithInvalidIdReturnsInvalidEntryID")

	t := conn.NewTransaction(uuid.New())

	accountID1 := "liability:clients:available:" + uuid.New().String()
	accountID2 := "liability:clients:available:" + uuid.New().String()

	invalidUUID := uuid.Nil
	t.AddEntry(uuid.New(), accountID1, vos.NewAccountVersion, vos.DebitOperation, 15000)
	t.AddEntry(invalidUUID, accountID2, vos.NewAccountVersion, vos.CreditOperation, 15000)

	err := conn.SaveTransaction(context.Background(), t)
	AssertTrue(ledger.ErrInvalidEntryID.Is(err))
}

func withoutEntriesReturnsInvalidEntriesNumber(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting withoutEntriesReturnsInvalidEntriesNumber")
	defer log.Println("finishing withoutEntriesReturnsInvalidEntriesNumber")

	t := conn.NewTransaction(uuid.New())

	err := conn.SaveTransaction(context.Background(), t)
	AssertTrue(ledger.ErrInvalidEntriesNumber.Is(err))
}

func entryWithInvalidVersion(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting EntryWithInvalidVersion")
	defer log.Println("finishing EntryWithInvalidVersion")

	t1 := conn.NewTransaction(uuid.New())

	accountID1 := "liability:clients:available:" + uuid.New().String()
	accountID2 := "liability:clients:available:" + uuid.New().String()

	t1.AddEntry(uuid.New(), accountID1, vos.NewAccountVersion, vos.DebitOperation, 15000)
	t1.AddEntry(uuid.New(), accountID2, vos.NewAccountVersion, vos.CreditOperation, 15000)
	err := conn.SaveTransaction(context.Background(), t1)

	t2 := conn.NewTransaction(uuid.New())
	invalidVersion := vos.NewAccountVersion
	t2.AddEntry(uuid.New(), accountID1, invalidVersion, vos.DebitOperation, 7500)
	t2.AddEntry(uuid.New(), accountID2, vos.AnyAccountVersion, vos.CreditOperation, 7500)
	err = conn.SaveTransaction(context.Background(), t2)

	AssertTrue(ledger.ErrInvalidVersion.Is(err))
}

func entryWithInvalidAccountStructure(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting entryWithInvalidAccountStructure")
	defer log.Println("finishing entryWithInvalidAccountStructure")

	t := conn.NewTransaction(uuid.New())

	invalidAccountID := "liability/clients/available/" + uuid.New().String()
	accountID2 := "liability:clients:available:" + uuid.New().String()

	t.AddEntry(uuid.New(), invalidAccountID, vos.NewAccountVersion, vos.DebitOperation, 15000)
	t.AddEntry(uuid.New(), accountID2, vos.NewAccountVersion, vos.CreditOperation, 15000)
	err := conn.SaveTransaction(context.Background(), t)

	AssertTrue(ledger.ErrInvalidAccountStructure.Is(err))
}

func transactionWithInvalidBalanceReturnsErrInvalidBalance(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting transactionWithInvalidBalance")
	defer log.Println("finishing transactionWithInvalidBalance")

	accountID1 := "liability:clients:available:" + uuid.New().String()
	accountID2 := "liability:clients:available:" + uuid.New().String()

	t := conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountID1, vos.NewAccountVersion, vos.DebitOperation, 15000)
	t.AddEntry(uuid.New(), accountID2, vos.NewAccountVersion, vos.CreditOperation, 25000)
	err := conn.SaveTransaction(context.Background(), t)
	AssertTrue(ledger.ErrInvalidBalance.Is(err))
}

func transactionWithInvalidIdempotencyKeyReturnsErrIdempotencyKeyViolation(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting transactionWithInvalidIdempotencyKeyReturnsErrIdempotencyKey")
	defer log.Println("finishing transactionWithInvalidIdempotencyKeyReturnsErrIdempotencyKey")

	accountID1 := "liability:clients:available:" + uuid.New().String()
	accountID2 := "liability:clients:available:" + uuid.New().String()

	idempotencyKey := uuid.New()

	t := conn.NewTransaction(uuid.New())
	t.AddEntry(idempotencyKey, accountID1, vos.NewAccountVersion, vos.DebitOperation, 15000)
	t.AddEntry(idempotencyKey, accountID2, vos.NewAccountVersion, vos.CreditOperation, 15000)
	err := conn.SaveTransaction(context.Background(), t)
	AssertTrue(ledger.ErrIdempotencyKeyViolation.Is(err))
}
