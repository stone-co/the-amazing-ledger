package querybuilder

import (
	"fmt"
	"strconv"
	"strings"
)

type QueryBuilder struct {
	queries map[int]string
	query   string
	numArgs int
}

func New(query string, numArgs int) QueryBuilder {
	qb := QueryBuilder{}
	qb.queries = make(map[int]string)
	qb.query = query
	qb.numArgs = numArgs

	return qb
}

func (q QueryBuilder) Init(numDefaultQueries int) {
	var offset = 2

	for i := offset; i < numDefaultQueries+offset; i++ {
		q.queries[i] = q.build(i)
	}
}

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
