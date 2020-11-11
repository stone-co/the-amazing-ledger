package entities

import (
	"testing"

	"gotest.tools/assert"
)

func TestAccountInfoBalance(t *testing.T) {
	accountName, _ := NewAccountName("liability:clients:available:user-1/block")

	accountInfo := NewAccountInfo(accountName.Name(), Version(3), 100, 50)

	assert.Equal(t, accountInfo.Balance(), 50)
}
