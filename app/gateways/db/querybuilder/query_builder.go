package querybuilder

import (
	"fmt"
	"strconv"
	"strings"
)

// QueryBuilder builds an insert query based on the number of arguments, how many values to insert,
// and an initial query defined as `insert into table_name(arg1, ..., argn) values %s;`.
// QueryBuilder maintains a map of queries that have already been built in memory to improve performance.
type QueryBuilder struct {
	queries map[int]string
	query   string
	numArgs int
}

// New creates a QueryBuilder and initializes it with the initial query, the number of arguments, and a 0-size queries map.
func New(query string, numArgs int) QueryBuilder {
	qb := QueryBuilder{}
	qb.queries = make(map[int]string)
	qb.query = query
	qb.numArgs = numArgs

	return qb
}

// Init initializes the queries map with n default queries.
func (q QueryBuilder) Init(n int) {
	var offset = 2

	for i := offset; i < n+offset; i++ {
		q.queries[i] = q.build(i)
	}
}

// Build builds a query based on the number of values to insert.
func (q QueryBuilder) Build(size int) string {
	query, ok := q.queries[size]
	if ok {
		return query
	}

	return q.build(size)
}

func (q QueryBuilder) build(size int) string {
	var sb strings.Builder

	sb.Grow(q.growLength(size))

	for i := 0; i < size; i++ {
		n := i * q.numArgs

		sb.WriteString("(")

		for j := 0; j < q.numArgs; j++ {
			sb.WriteString("$")
			sb.WriteString(strconv.Itoa(n + j + 1))

			if j != q.numArgs-1 {
				sb.WriteString(", ")
			}
		}

		if i != size-1 {
			sb.WriteString("),\n")
		}
	}

	sb.WriteString(")")

	query := fmt.Sprintf(q.query, sb.String())
	q.queries[size] = query

	return query
}

func (q QueryBuilder) growLength(size int) int {
	digitsLength := q.numArgs * 3
	numberOfDollars := q.numArgs
	numberOfCommasAndSpaces := (q.numArgs - 1) * 2
	numberOfParenthesis := 2

	// valuesLength is the length of a value to be inserted ($1, ..., $2)
	valuesLength := digitsLength + numberOfDollars + numberOfCommasAndSpaces + numberOfParenthesis

	queryLength := len(q.query) + valuesLength*size + (size - 1) + (size-1)*2 - 2

	return queryLength
}
