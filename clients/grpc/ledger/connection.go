package ledger

import (
	"context"
	"fmt"

	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"

	"google.golang.org/grpc"
)

type Connection struct {
	conn   *grpc.ClientConn
	client proto.LedgerServiceClient
}

func Connect(ctx context.Context, host string, port int) (*Connection, error) {
	target := fmt.Sprintf("%s:%d", host, port)

	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnectionFailed, err)
	}

	client := proto.NewLedgerServiceClient(conn)

	return &Connection{
		conn:   conn,
		client: client,
	}, nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
