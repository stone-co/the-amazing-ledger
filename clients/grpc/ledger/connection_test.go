package ledger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectFailed(t *testing.T) {
	host := "localhost"
	port := 4000

	_, err := Connect(host, port)
	fmt.Println(err)

	assert.True(t, ErrConnectionFailed.Is(err))
}
