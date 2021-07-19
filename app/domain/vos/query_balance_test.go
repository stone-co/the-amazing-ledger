package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryBalance(t *testing.T) {
	query := "liability.clients.available.user_1.*"
	accountQuery, err := NewAccount(query)
	assert.NoError(t, err)

	queryBalance := NewQueryBalance(accountQuery, 100)
	assert.Equal(t, QueryBalance{accountQuery, 100}, queryBalance)
}
