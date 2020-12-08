package ledger

import (
	"fmt"
	"time"

	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"

	"google.golang.org/grpc"
)

type Connection struct {
	conn   *grpc.ClientConn
	client proto.LedgerServiceClient
}

func Connect(host string, port int) (*Connection, error) {

	target := fmt.Sprintf("%s:%d", host, port)

	// TODO: uses DialContext instead Dial
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*5))
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
