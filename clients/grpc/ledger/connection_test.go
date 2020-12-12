package ledger

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectFailed(t *testing.T) {
	ctx := context.Background()
	host := "localhost"
	port := 4000

	_, err := Connect(ctx, host, port)
	assert.True(t, errors.Is(err, ErrConnectionFailed))
}
