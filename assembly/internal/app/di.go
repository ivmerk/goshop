package app

import (
	"context"

	"github.com/IBM/sarama"

	"github.com/ivmerk/goshop/assembly/internal/config"
	"github.com/ivmerk/goshop/assembly/internal/service"
	consumerSvc "github.com/ivmerk/goshop/assembly/internal/service/consumer/order_paid_consumer"
	producerSvc "github.com/ivmerk/goshop/assembly/internal/service/producer/order_assembled_producer"
	"github.com/ivmerk/goshop/platform/pkg/closer"
	"github.com/ivmerk/goshop/platform/pkg/kafka"
	kafkaConsumer "github.com/ivmerk/goshop/platform/pkg/kafka/consumer"
	kafkaProducer "github.com/ivmerk/goshop/platform/pkg/kafka/producer"
	"github.com/ivmerk/goshop/platform/pkg/logger"
	kafkaMiddleware "github.com/ivmerk/goshop/platform/pkg/middleware/kafka"
)

type diContainer struct {
	consumerGroup sarama.ConsumerGroup
	syncProducer  sarama.SyncProducer

	orderPaidConsumer kafka.Consumer
	producerSvc       service.ProducerService
	consumerSvc       service.ConsumerService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) ConsumerGroup() (sarama.ConsumerGroup, error) {
	if d.consumerGroup == nil {
		cfg := config.AppConfig()
		g, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers(), cfg.OrderPaidConsumer.GroupID(), cfg.OrderPaidConsumer.Config())
		if err != nil {
			return nil, err
		}
		closer.AddNamed("Kafka consumer group", func(context.Context) error { return g.Close() })
		d.consumerGroup = g
	}
	return d.consumerGroup, nil
}

func (d *diContainer) SyncProducer() (sarama.SyncProducer, error) {
	if d.syncProducer == nil {
		cfg := config.AppConfig()
		p, err := sarama.NewSyncProducer(cfg.Kafka.Brokers(), cfg.OrderAssembledProducer.Config())
		if err != nil {
			return nil, err
		}
		closer.AddNamed("Kafka sync producer", func(context.Context) error { return p.Close() })
		d.syncProducer = p
	}
	return d.syncProducer, nil
}

func (d *diContainer) ProducerService() (service.ProducerService, error) {
	if d.producerSvc == nil {
		sp, err := d.SyncProducer()
		if err != nil {
			return nil, err
		}
		kp := kafkaProducer.NewProducer(sp, config.AppConfig().OrderAssembledProducer.Topic(), logger.Logger())
		d.producerSvc = producerSvc.NewService(kp)
	}
	return d.producerSvc, nil
}

func (d *diContainer) OrderPaidConsumer() (kafka.Consumer, error) {
	if d.orderPaidConsumer == nil {
		g, err := d.ConsumerGroup()
		if err != nil {
			return nil, err
		}
		d.orderPaidConsumer = kafkaConsumer.NewConsumer(
			g,
			[]string{config.AppConfig().OrderPaidConsumer.Topic()},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.orderPaidConsumer, nil
}

func (d *diContainer) ConsumerService() (service.ConsumerService, error) {
	if d.consumerSvc == nil {
		c, err := d.OrderPaidConsumer()
		if err != nil {
			return nil, err
		}
		p, err := d.ProducerService()
		if err != nil {
			return nil, err
		}
		d.consumerSvc = consumerSvc.NewService(c, p)
	}
	return d.consumerSvc, nil
}
