package env

import (
	"os"

	"github.com/pkg/errors"
)

const (
	orderPaidProducerTopicEnvName = "ORDER_PAID_PRODUCER_TOPIC"
)

type orderPaidProducerConfig struct {
	topic string
}

func NewOrderPaidProducerConfig() (*orderPaidProducerConfig, error) {
	topic := os.Getenv(orderPaidProducerTopicEnvName)
	if len(topic) == 0 {
		return nil, errors.New("order paid producer topic is required")
	}

	return &orderPaidProducerConfig{
		topic: topic,
	}, nil
}

func (cfg *orderPaidProducerConfig) Topic() string {
	return cfg.topic
}
