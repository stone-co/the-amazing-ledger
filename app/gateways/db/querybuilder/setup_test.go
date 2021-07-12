package querybuilder

import (
	"os"
	"testing"
)

var (
	numArgs int
	query   string
)

func TestMain(m *testing.M) {
	numArgs = 10
	query = `
insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company, metadata)
values %s;`

	os.Exit(m.Run())
}
