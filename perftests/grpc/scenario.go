package main

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// We are using GHZ (https://ghz.sh/) to run some perf tests with the Ledger.

// Scenario:
// - we create a "stone" fake account (assets:aaa:bbb:sa)
// - we create 1k fake client accounts (liability:aaa:bbb:UUID)

// For each client account above, we create a transaction that transfers money from the client account to the stone account.
// These transactions are executed N times.

type Scenario struct {
	stoneAccount   string
	clientAccounts []string
	transactions   string
}

func NewScenario(totalClientAccounts, totalRequests int) *Scenario {
	s := &Scenario{}
	s.setup(totalClientAccounts, totalRequests)
	s.defineSteps()
	return s
}

func (s *Scenario) setup(totalClientAccounts, totalRequests int) {
	// Create the Stone account.
	s.stoneAccount = "assets:aaa:bbb:sa"

	// Create N client accounts.
	s.clientAccounts = make([]string, totalClientAccounts)
	for i := 0; i < totalClientAccounts; i++ {
		s.clientAccounts[i] = "liability:aaa:bbb:" + uuid.New().String()
	}
}

func (s *Scenario) defineSteps() {
	var series strings.Builder
	series.WriteString("[")

	// For each client account above, execute:
	for _, ca := range s.clientAccounts {
		// TODO: 1 get account balance on the client account

		// 1 transaction that transfer money from CA to SA
		series.WriteString(s.transferTransaction(ca) + ",")
	}

	s.transactions = series.String()
	s.transactions = s.transactions[:len(s.transactions)-1] + "]"
}

func (s *Scenario) transferTransaction(toAccountID string) string {
	amount := 20000
	e1 := entryAsString("{{newUUID}}", toAccountID, 0, "OPERATION_DEBIT", amount)
	e2 := entryAsString("{{newUUID}}", s.stoneAccount, 0, "OPERATION_CREDIT", amount)
	tr := fmt.Sprintf(`{"id":"{{newUUID}}", "entries":[%s,%s]}`, e1, e2)

	return tr
}

func entryAsString(id string, accountID string, expectedVersion int, operation string, amount int) string {
	return fmt.Sprintf(`{"id":"%s","account_id":"%s", "expected_version": %d, "operation": "%s", "amount": %d}`, id, accountID, expectedVersion, operation, amount)
}

func (s Scenario) JSON() string {
	return s.transactions
}
