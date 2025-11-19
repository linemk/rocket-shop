package order_consumer

import (
	"context"

	"go.uber.org/zap"

	platformKafka "github.com/linemk/rocket-shop/platform/pkg/kafka"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type HandlerFunc func(ctx context.Context, msg platformKafka.Message) error

type consumer struct {
	kafkaConsumer platformKafka.Consumer
	handler       HandlerFunc
	logger        Logger
}

func NewConsumer(
	kafkaConsumer platformKafka.Consumer,
	handler HandlerFunc,
	logger Logger,
) *consumer {
	return &consumer{
		kafkaConsumer: kafkaConsumer,
		handler:       handler,
		logger:        logger,
	}
}

func (c *consumer) RunConsumers(ctx context.Context) error {
	c.logger.Info(ctx, "Starting Kafka consumer for ShipAssembled events")

	if err := c.kafkaConsumer.Consume(ctx, platformKafka.MessageHandler(c.handler)); err != nil {
		c.logger.Error(ctx, "Kafka consumer error", zap.Error(err))
		return err
	}

	return nil
}
