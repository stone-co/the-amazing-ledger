package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryBalance(t *testing.T) {
	query := "liability.clients.available.user-1.*"
	accountQuery, err := NewAccountQuery(query)
	assert.NoError(t, err)

	queryBalance := NewQueryBalance(accountQuery, 100)
	assert.Equal(t, QueryBalance{accountQuery, 100}, queryBalance)
}
