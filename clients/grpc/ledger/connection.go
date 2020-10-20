package ledger

import (
	"fmt"

	"google.golang.org/grpc"

	pb "github.com/stone-co/the-amazing-ledger/pkg/gateways/grpc/proto/ledger"
)

type Connection struct {
	conn   *grpc.ClientConn
	client pb.LedgerServiceClient
}

func Connect(host string, port int) (*Connection, error) {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %s", err)
	}

	// TODO: is's ok? 1 for all calls?
	client := pb.NewLedgerServiceClient(conn)

	return &Connection{
		conn:   conn,
		client: client,
	}, nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
