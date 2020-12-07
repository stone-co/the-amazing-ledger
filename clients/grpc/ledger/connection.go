package ledger

import (
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"

	"google.golang.org/grpc"
)

type Connection struct {
	conn   *grpc.ClientConn
	client proto.LedgerServiceClient
}

func Connect(host string, port int) (*Connection, error) {
	conn, err := grpc.Dial(":3000", grpc.WithInsecure())
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
