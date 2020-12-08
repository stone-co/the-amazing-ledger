package ledger

import (
	"context"
	"fmt"
	"time"

	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"

	"google.golang.org/grpc"
)

type Connection struct {
	conn   *grpc.ClientConn
	client proto.LedgerServiceClient
}

func Connect(ctx context.Context, host string, port int) (*Connection, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	target := fmt.Sprintf("%s:%d", host, port)

	// TODO: uses DialContext instead Dial
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, ErrConnectionFailed.cause(err)
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
