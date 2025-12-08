package grpcserver

import (
	"fmt"
	"net"
)

// NewListener создает новый слушатель для gRPC сервера
func NewListener(addr string) (net.Listener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	return listener, nil
}
