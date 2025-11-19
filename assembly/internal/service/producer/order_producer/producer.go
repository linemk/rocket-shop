package order_producer

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/linemk/rocket-shop/assembly/internal/converter/kafka"
	"github.com/linemk/rocket-shop/assembly/internal/model"
	platformKafka "github.com/linemk/rocket-shop/platform/pkg/kafka"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type producer struct {
	kafkaProducer platformKafka.Producer
	logger        Logger
}

func NewProducer(kafkaProducer platformKafka.Producer, logger Logger) *producer {
	return &producer{
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (p *producer) SendShipAssembled(ctx context.Context, event *model.ShipAssembledEvent) error {
	data, err := kafka.EncodeShipAssembled(event)
	if err != nil {
		p.logger.Error(ctx, "Failed to encode ShipAssembled event", zap.Error(err))
		return err
	}

	key := []byte(uuid.New().String())
	if err := p.kafkaProducer.Send(ctx, key, data); err != nil {
		p.logger.Error(ctx, "Failed to send ShipAssembled event to Kafka", zap.Error(err))
		return err
	}

	p.logger.Info(ctx, "ShipAssembled event sent successfully",
		zap.String("order_uuid", event.OrderUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	return nil
}
