package env

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

const (
	kafkaBrokersEnvName = "KAFKA_BROKERS"
)

type kafkaConfig struct {
	brokers []string
}

func NewKafkaConfig() (*kafkaConfig, error) {
	brokersStr := os.Getenv(kafkaBrokersEnvName)
	if len(brokersStr) == 0 {
		return nil, errors.New("kafka brokers is required")
	}

	brokers := strings.Split(brokersStr, ",")

	return &kafkaConfig{
		brokers: brokers,
	}, nil
}

func (cfg *kafkaConfig) Brokers() []string {
	return cfg.brokers
}
