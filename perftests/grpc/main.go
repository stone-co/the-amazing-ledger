package main

import (
	"fmt"
	"os"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
)

func main() {

	totalClientAccounts := 1000
	totalRequests := 10000
	concurrency := 20

	scenario := NewScenario(totalClientAccounts, totalRequests)

	report, err := runner.Run(
		"ledger.LedgerService.CreateTransaction",
		"0.0.0.0:3000",
		runner.WithProtoFile("../../proto/ledger/ledger.proto", []string{"../../third_party/googleapis/"}),
		runner.WithDataFromJSON(scenario.JSON()),
		// runner.WithDataFromReader(scenario),
		runner.WithInsecure(true),
		runner.WithTotalRequests(uint(totalRequests)),
		runner.WithConcurrency(uint(concurrency)),
		// runner.WithTimeout(0),
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	printer := printer.ReportPrinter{
		Out:    os.Stdout,
		Report: report,
	}

	printer.Print("summary")
}
