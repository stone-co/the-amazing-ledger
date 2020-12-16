package ledger

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnectFailed(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	host := "localhost"
	port := 4000

	_, err := Connect(ctx, host, port)
	assert.True(t, errors.Is(err, ErrConnectionFailed))
}
