package vos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccountBalance(t *testing.T) {
	accountName, _ := NewAccountPath("liability.clients.available.user-1.block")

	accountBalance := NewAccountBalance(accountName, Version(3), 100, 50, time.Now())

	assert.Equal(t, 50, accountBalance.Balance())
}
