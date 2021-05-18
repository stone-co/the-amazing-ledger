package rpc

import (
	"os"
	"testing"
)

const ValidTransactionID = "35b2ac20-d733-461c-810a-d7f4acf6316f"
const ValidAccountID = "liability.clients.available.96a131a8-c4ac-495e-8971-fcecdbdd003a"

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
