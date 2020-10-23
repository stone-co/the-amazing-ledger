package ledger

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto"
)

type Connection struct {
	conn   *grpc.ClientConn
	client proto.LedgerServiceClient
}

func Connect(host string, port int) (*Connection, error) {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %s", err)
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
