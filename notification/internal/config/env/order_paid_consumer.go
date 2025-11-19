package env

import (
	"os"

	"github.com/pkg/errors"
)

const (
	orderPaidConsumerTopicEnvName   = "ORDER_PAID_CONSUMER_TOPIC"
	orderPaidConsumerGroupIDEnvName = "ORDER_PAID_CONSUMER_GROUP_ID"
)

type orderPaidConsumerConfig struct {
	topic   string
	groupID string
}

func NewOrderPaidConsumerConfig() (*orderPaidConsumerConfig, error) {
	topic := os.Getenv(orderPaidConsumerTopicEnvName)
	if len(topic) == 0 {
		return nil, errors.New("order paid consumer topic is required")
	}

	groupID := os.Getenv(orderPaidConsumerGroupIDEnvName)
	if len(groupID) == 0 {
		return nil, errors.New("order paid consumer group id is required")
	}

	return &orderPaidConsumerConfig{
		topic:   topic,
		groupID: groupID,
	}, nil
}

func (cfg *orderPaidConsumerConfig) Topic() string {
	return cfg.topic
}

func (cfg *orderPaidConsumerConfig) GroupID() string {
	return cfg.groupID
}
