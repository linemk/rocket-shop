package iamclient

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/linemk/rocket-shop/shared/pkg/proto/auth/v1"
	userv1 "github.com/linemk/rocket-shop/shared/pkg/proto/user/v1"
)

type Client struct {
	conn    *grpc.ClientConn
	authSvc authv1.AuthServiceClient
	userSvc userv1.UserServiceClient
}

func New(ctx context.Context, addr string) (*Client, error) {
	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		authSvc: authv1.NewAuthServiceClient(conn),
		userSvc: userv1.NewUserServiceClient(conn),
	}, nil
}

func (c *Client) Auth() authv1.AuthServiceClient {
	return c.authSvc
}

func (c *Client) User() userv1.UserServiceClient {
	return c.userSvc
}

func (c *Client) Close() error {
	return c.conn.Close()
}
