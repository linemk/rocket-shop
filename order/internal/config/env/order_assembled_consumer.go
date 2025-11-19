package env

import (
	"os"

	"github.com/pkg/errors"
)

const (
	orderAssembledConsumerTopicEnvName   = "ORDER_ASSEMBLED_CONSUMER_TOPIC"
	orderAssembledConsumerGroupIDEnvName = "ORDER_ASSEMBLED_CONSUMER_GROUP_ID"
)

type orderAssembledConsumerConfig struct {
	topic   string
	groupID string
}

func NewOrderAssembledConsumerConfig() (*orderAssembledConsumerConfig, error) {
	topic := os.Getenv(orderAssembledConsumerTopicEnvName)
	if len(topic) == 0 {
		return nil, errors.New("order assembled consumer topic is required")
	}

	groupID := os.Getenv(orderAssembledConsumerGroupIDEnvName)
	if len(groupID) == 0 {
		return nil, errors.New("order assembled consumer group id is required")
	}

	return &orderAssembledConsumerConfig{
		topic:   topic,
		groupID: groupID,
	}, nil
}

func (cfg *orderAssembledConsumerConfig) Topic() string {
	return cfg.topic
}

func (cfg *orderAssembledConsumerConfig) GroupID() string {
	return cfg.groupID
}
