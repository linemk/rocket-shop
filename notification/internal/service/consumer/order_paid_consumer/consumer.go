package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	platformKafka "github.com/linemk/rocket-shop/platform/pkg/kafka"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Consumer struct {
	kafkaConsumer platformKafka.Consumer
	handler       platformKafka.MessageHandler
	logger        Logger
}

func NewConsumer(
	kafkaConsumer platformKafka.Consumer,
	handler platformKafka.MessageHandler,
	logger Logger,
) *Consumer {
	return &Consumer{
		kafkaConsumer: kafkaConsumer,
		handler:       handler,
		logger:        logger,
	}
}

func (c *Consumer) RunConsumer(ctx context.Context) error {
	c.logger.Info(ctx, "Starting Kafka consumer for OrderPaid events")

	if err := c.kafkaConsumer.Consume(ctx, c.handler); err != nil {
		c.logger.Error(ctx, "Kafka consumer error", zap.Error(err))
		return err
	}

	return nil
}
