package rpc

import (
	"os"
	"testing"
)

const ValidTransactionID = "35b2ac20-d733-461c-810a-d7f4acf6316f"
const ValidAccountID = "liability.clients.available.96a131a8_c4ac_495e_8971_fcecdbdd003a"

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
