package v1

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventory_v1 "github.com/linemk/rocket-shop/shared/pkg/proto/inventory/v1"
)

type Client struct {
	client inventory_v1.InventoryServiceClient
	conn   *grpc.ClientConn // нужен для закрытия соединения
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to InventoryService: %w", err)
	}

	client := inventory_v1.NewInventoryServiceClient(conn)

	return &Client{
		client: client,
		conn:   conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// Проверяем, что Client реализует интерфейс InventoryClient
var _ InventoryClient = (*Client)(nil)
