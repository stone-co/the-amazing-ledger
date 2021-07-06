package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/clients/grpc/ledger"
)

func getSyntheticReportFullPath(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting GetSyntheticReport")
	defer log.Println("finishing GetSyntheticReport")

	// expectedBalance := 1000
	accountPathBase := "liability.stone.clients"
	accountPathOne := accountPathBase + "." + strings.ReplaceAll(uuid.New().String(), "-", "")
	accountPathTwo := accountPathBase + "." + strings.ReplaceAll(uuid.New().String(), "-", "")

	// Define a new transaction with 2 entries
	t := conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, vos.NextAccountVersion, vos.CreditOperation, 1000)
	t.AddEntry(uuid.New(), accountPathTwo, vos.NextAccountVersion, vos.DebitOperation, 1000)
	t.Message.Event = 1 // generic event
	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	now := time.Now().UnixNano()
	report, err := conn.GetSyntheticReport(context.Background(), accountPathBase, now, now)
	fmt.Printf("> report: %v\n\n", report)
	AssertTrue(report != nil)

	paths := report.Paths()
	AssertTrue(paths != nil)

	AssertEqual(accountPathOne, paths[0].Account)
	AssertEqual(int64(1000), paths[0].Credit)
	AssertEqual(int64(0), paths[0].Debit)
}

func getSyntheticReportSubgroup(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting GetSyntheticReport Subgroup")
	defer log.Println("finishing GetSyntheticReport Subgroup")

	// expectedBalance := 1000
	accountPathOne := "liability.stone.example." + strings.ReplaceAll(uuid.New().String(), "-", "")
	accountPathTwo := "liability.stone.example." + strings.ReplaceAll(uuid.New().String(), "-", "")

	// Define a new transaction with 2 entries
	t := conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, vos.NextAccountVersion, vos.CreditOperation, 1000)
	t.AddEntry(uuid.New(), accountPathTwo, vos.NextAccountVersion, vos.DebitOperation, 1000)
	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	now := time.Now().UnixNano()
	report, err := conn.GetSyntheticReport(context.Background(), accountPathOne, now, now)
	fmt.Printf("> report: %v\n\n", report)
	AssertTrue(report != nil)

	paths := report.Paths()
	AssertTrue(paths != nil)

	AssertEqual(accountPathOne, paths[0].Account)
	AssertEqual(int64(1000), paths[0].Credit)
	AssertEqual(int64(0), paths[0].Debit)
}

func getSyntheticReportGroup(log *logrus.Entry, conn *ledger.Connection) {
	log.Println("starting GetSyntheticReport Subgroup")
	defer log.Println("finishing GetSyntheticReport Subgroup")

	// expectedBalance := 1000
	accountPathOne := "liability:xpto:clients:" + uuid.New().String()
	accountPathTwo := "liability:xpto:clients:" + uuid.New().String()

	// Define a new transaction with 2 entries
	t := conn.NewTransaction(uuid.New())
	t.AddEntry(uuid.New(), accountPathOne, vos.NextAccountVersion, vos.CreditOperation, 1000)
	t.AddEntry(uuid.New(), accountPathTwo, vos.NextAccountVersion, vos.DebitOperation, 1000)
	err := conn.SaveTransaction(context.Background(), t)
	AssertEqual(nil, err)

	now := time.Now().UnixNano()
	report, err := conn.GetSyntheticReport(context.Background(), accountPathOne, now, now)
	fmt.Printf("> report one: %v\n\n", report)
	AssertTrue(report != nil)

	paths := report.Paths()
	AssertTrue(paths != nil)

	AssertEqual(accountPathOne, paths[0].Account)
	AssertEqual(int64(1000), paths[0].Credit)
	AssertEqual(int64(0), paths[0].Debit)

	reportTwo, err := conn.GetSyntheticReport(context.Background(), accountPathTwo, now, now)
	fmt.Printf("> report two: %v\n\n", reportTwo)
	AssertTrue(reportTwo != nil)

	pathsTwo := reportTwo.Paths()
	AssertTrue(pathsTwo != nil)

	AssertEqual(accountPathTwo, pathsTwo[0].Account)
	AssertEqual(int64(0), pathsTwo[0].Credit)
	AssertEqual(int64(1000), pathsTwo[0].Debit)

}
