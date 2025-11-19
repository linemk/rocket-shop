package env

import (
	"os"

	"github.com/pkg/errors"
)

const (
	orderAssembledProducerTopicEnvName = "ORDER_ASSEMBLED_PRODUCER_TOPIC"
)

type orderAssembledProducerConfig struct {
	topic string
}

func NewOrderAssembledProducerConfig() (*orderAssembledProducerConfig, error) {
	topic := os.Getenv(orderAssembledProducerTopicEnvName)
	if len(topic) == 0 {
		return nil, errors.New("order assembled producer topic is required")
	}

	return &orderAssembledProducerConfig{
		topic: topic,
	}, nil
}

func (cfg *orderAssembledProducerConfig) Topic() string {
	return cfg.topic
}
