package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountBalance(t *testing.T) {
	accountName, _ := NewAccountName("liability:clients:available:user-1/block")

	accountBalance := NewAccountBalance(*accountName, Version(3), 100, 50)

	assert.Equal(t, 50, accountBalance.Balance())
}
