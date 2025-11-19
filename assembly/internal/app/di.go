package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/linemk/rocket-shop/assembly/internal/config"
	"github.com/linemk/rocket-shop/assembly/internal/service"
	"github.com/linemk/rocket-shop/assembly/internal/service/consumer/order_consumer"
	"github.com/linemk/rocket-shop/assembly/internal/service/producer/order_producer"
	"github.com/linemk/rocket-shop/platform/pkg/closer"
	"github.com/linemk/rocket-shop/platform/pkg/kafka/consumer"
	"github.com/linemk/rocket-shop/platform/pkg/kafka/producer"
	"github.com/linemk/rocket-shop/platform/pkg/logger"
	kafkaMiddleware "github.com/linemk/rocket-shop/platform/pkg/middleware/kafka"
)

type diContainer struct {
	consumerService      service.ConsumerService
	orderProducerService service.OrderProducerService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) ConsumerService(ctx context.Context) service.ConsumerService {
	if d.consumerService == nil {
		// Создаем Kafka consumer group
		saramaConfig := sarama.NewConfig()
		saramaConfig.Version = sarama.V2_6_0_0
		saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			saramaConfig,
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create Kafka consumer group: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return consumerGroup.Close()
		})

		// Создаем Kafka consumer с middleware
		kafkaConsumer := consumer.NewConsumer(
			consumerGroup,
			[]string{config.AppConfig().OrderPaidConsumer.Topic()},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)

		// Создаем handler
		handler := order_consumer.NewHandler(d.OrderProducerService(ctx), logger.Logger())

		d.consumerService = order_consumer.NewConsumer(kafkaConsumer, handler, logger.Logger())
	}

	return d.consumerService
}

func (d *diContainer) OrderProducerService(ctx context.Context) service.OrderProducerService {
	if d.orderProducerService == nil {
		// Создаем Kafka sync producer
		saramaConfig := sarama.NewConfig()
		saramaConfig.Version = sarama.V2_6_0_0
		saramaConfig.Producer.Return.Successes = true
		saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
		saramaConfig.Producer.Retry.Max = 5

		syncProducer, err := sarama.NewSyncProducer(config.AppConfig().Kafka.Brokers(), saramaConfig)
		if err != nil {
			panic(fmt.Sprintf("failed to create Kafka sync producer: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return syncProducer.Close()
		})

		kafkaProducer := producer.NewProducer(
			syncProducer,
			config.AppConfig().OrderAssembledProducer.Topic(),
			logger.Logger(),
		)

		d.orderProducerService = order_producer.NewProducer(kafkaProducer, logger.Logger())
	}

	return d.orderProducerService
}
