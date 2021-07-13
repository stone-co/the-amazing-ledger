package postgres

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stone-co/the-amazing-ledger/app/domain/entities"
	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	"github.com/stone-co/the-amazing-ledger/app/tests"
	"github.com/stretchr/testify/assert"
)

var pgDocker *tests.PostgresDocker

func TestMain(m *testing.M) {
	pgDocker = tests.SetupTest("./migrations")

	_, err := pgDocker.DB.Exec(context.Background(), `insert into event (id, name) values (1, 'default');`)
	if err != nil {
		log.Fatalf("could not insert default event values: %v", err)
	}

	exitCode := m.Run()

	tests.RemoveContainer(pgDocker)

	os.Exit(exitCode)
}

func createEntry(t *testing.T, operation vos.OperationType, account string, version vos.Version) entities.Entry {
	entry, err := entities.NewEntry(
		uuid.New(),
		operation,
		account,
		version,
		100,
		json.RawMessage(`{}`),
	)
	assert.NoError(t, err)

	return entry
}

func createTransaction(t *testing.T, ctx context.Context, r *LedgerRepository, entries ...entities.Entry) error {
	tx, err := entities.NewTransaction(uuid.New(), uint32(1), "abc", time.Now(), entries...)
	assert.NoError(t, err)

	return r.CreateTransaction(ctx, tx)
}
