package ledger

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectFailed(t *testing.T) {
	ctx := context.Background()
	host := "localhost"
	port := 4000

	_, err := Connect(ctx, host, port)
	fmt.Println(err)

	assert.True(t, ErrConnectionFailed.Is(err))
}
