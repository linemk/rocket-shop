package v1

import (
	"fmt"

	payment_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client payment_v1.PaymentServiceClient
	conn   *grpc.ClientConn
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PaymentService: %w", err)
	}

	client := payment_v1.NewPaymentServiceClient(conn)

	return &Client{
		client: client,
		conn:   conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
