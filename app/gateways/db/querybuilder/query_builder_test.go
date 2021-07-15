package querybuilder

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	qb := New(query, numArgs)

	assert.Len(t, qb.queries, 0)
	assert.Equal(t, qb.query, query)
	assert.Equal(t, qb.numArgs, numArgs)
}

func TestQueryBuilder_Init(t *testing.T) {
	testCases := []struct {
		name              string
		numDefaultQueries int
		expectedLen       int
	}{
		{
			name:              "should create query builder with 1 default query",
			numDefaultQueries: 1,
			expectedLen:       1,
		},
		{
			name:              "should create query builder with 2 default queries",
			numDefaultQueries: 2,
			expectedLen:       2,
		},
		{
			name:              "should create query builder with 3 default queries",
			numDefaultQueries: 3,
			expectedLen:       3,
		},
		{
			name:              "should create query builder with 4 default queries",
			numDefaultQueries: 4,
			expectedLen:       4,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			qb := New(query, numArgs)
			qb.Init(tt.numDefaultQueries)
			assert.Len(t, qb.queries, tt.expectedLen)
		})
	}
}

func TestQueryBuilder_Build(t *testing.T) {
	t.Run("should return a query if not cached successfully", func(t *testing.T) {
		qb := New(query, numArgs)
		assert.Len(t, qb.queries, 0)

		qb.Build(2)
		assert.Len(t, qb.queries, 1)
	})

	t.Run("should return a query if it's cached successfully", func(t *testing.T) {
		qb := New(query, numArgs)
		qb.Init(1)
		assert.Len(t, qb.queries, 1)

		qb.Build(2)
		assert.Len(t, qb.queries, 1)
	})
}

func TestQueryBuilder_build(t *testing.T) {
	testCases := []struct {
		name     string
		size     int
		expected string
	}{
		{
			name: "should create query with 2 entries successfully",
			size: 2,
			expected: `
				insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company, metadata)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10),
				($11, $12, $13, $14, $15, $16, $17, $18, $19, $20);`,
		},
		{
			name: "should create query with 3 entries successfully",
			size: 3,
			expected: `
				insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company, metadata)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10),
				($11, $12, $13, $14, $15, $16, $17, $18, $19, $20),
				($21, $22, $23, $24, $25, $26, $27, $28, $29, $30);`,
		},
		{
			name: "should create query with 4 entries successfully",
			size: 4,
			expected: `
				insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company, metadata)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10),
				($11, $12, $13, $14, $15, $16, $17, $18, $19, $20),
				($21, $22, $23, $24, $25, $26, $27, $28, $29, $30),
				($31, $32, $33, $34, $35, $36, $37, $38, $39, $40);`,
		},
		{
			name: "should create query with 5 entries successfully",
			size: 5,
			expected: `
				insert into entry (id, tx_id, event, operation, version, amount, competence_date, account, company, metadata)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10),
				($11, $12, $13, $14, $15, $16, $17, $18, $19, $20),
				($21, $22, $23, $24, $25, $26, $27, $28, $29, $30),
				($31, $32, $33, $34, $35, $36, $37, $38, $39, $40),
				($41, $42, $43, $44, $45, $46, $47, $48, $49, $50);`,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			qb := New(query, numArgs)
			got := qb.build(tt.size)
			want := strings.ReplaceAll(tt.expected, "\t", "")

			assert.Equal(t, want, got)
		})
	}
}
