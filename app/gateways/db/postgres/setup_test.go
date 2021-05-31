package postgres

import (
	"os"
	"testing"

	"github.com/stone-co/the-amazing-ledger/app/tests"
)

var pgDocker *tests.PostgresDocker

func TestMain(m *testing.M) {
	pgDocker = tests.SetupTest("./migrations")

	exitCode := m.Run()

	tests.RemoveContainer(pgDocker)

	os.Exit(exitCode)
}
