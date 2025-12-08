package env

import "os"

const iamGRPCAddressEnv = "IAM_GRPC_ADDRESS"

type iamGRPCConfig struct {
	address string
}

func NewIAMGRPCConfig() (*iamGRPCConfig, error) {
	address := os.Getenv(iamGRPCAddressEnv)
	if address == "" {
		address = "localhost:50051"
	}
	return &iamGRPCConfig{
		address: address,
	}, nil
}

func (c *iamGRPCConfig) Address() string {
	return c.address
}
