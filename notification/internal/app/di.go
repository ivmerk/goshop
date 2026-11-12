package app

import (
	"context"

	"github.com/IBM/sarama"

	httpClient "github.com/ivmerk/goshop/notification/internal/client/http"
	telegramClient "github.com/ivmerk/goshop/notification/internal/client/http/telegram"
	"github.com/ivmerk/goshop/notification/internal/config"
	"github.com/ivmerk/goshop/notification/internal/service"
	orderAssembledConsumerSvc "github.com/ivmerk/goshop/notification/internal/service/consumer/order_assembled_consumer"
	orderPaidConsumerSvc "github.com/ivmerk/goshop/notification/internal/service/consumer/order_paid_consumer"
	telegramSvc "github.com/ivmerk/goshop/notification/internal/service/telegram"
	"github.com/ivmerk/goshop/platform/pkg/closer"
	"github.com/ivmerk/goshop/platform/pkg/kafka"
	kafkaConsumer "github.com/ivmerk/goshop/platform/pkg/kafka/consumer"
	"github.com/ivmerk/goshop/platform/pkg/logger"
	kafkaMiddleware "github.com/ivmerk/goshop/platform/pkg/middleware/kafka"
)

type diContainer struct {
	telegramClient httpClient.TelegramClient
	telegramSvc    service.TelegramService

	orderPaidConsumerGroup sarama.ConsumerGroup
	orderPaidConsumer      kafka.Consumer
	orderPaidConsumerSvc   service.OrderPaidConsumerService

	orderAssembledConsumerGroup sarama.ConsumerGroup
	orderAssembledConsumer      kafka.Consumer
	orderAssembledConsumerSvc   service.OrderAssembledConsumerService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) TelegramClient() (httpClient.TelegramClient, error) {
	if d.telegramClient == nil {
		c, err := telegramClient.NewClient(config.AppConfig().TelegramBot.Token())
		if err != nil {
			return nil, err
		}
		d.telegramClient = c
	}
	return d.telegramClient, nil
}

func (d *diContainer) TelegramService() (service.TelegramService, error) {
	if d.telegramSvc == nil {
		c, err := d.TelegramClient()
		if err != nil {
			return nil, err
		}
		s, err := telegramSvc.NewService(c, config.AppConfig().TelegramBot.ChatID())
		if err != nil {
			return nil, err
		}
		d.telegramSvc = s
	}
	return d.telegramSvc, nil
}

func (d *diContainer) OrderPaidConsumerGroup() (sarama.ConsumerGroup, error) {
	if d.orderPaidConsumerGroup == nil {
		cfg := config.AppConfig()
		g, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers(), cfg.OrderPaidConsumer.GroupID(), cfg.OrderPaidConsumer.Config())
		if err != nil {
			return nil, err
		}
		closer.AddNamed("Kafka consumer group (order paid)", func(context.Context) error { return g.Close() })
		d.orderPaidConsumerGroup = g
	}
	return d.orderPaidConsumerGroup, nil
}

func (d *diContainer) OrderPaidConsumer() (kafka.Consumer, error) {
	if d.orderPaidConsumer == nil {
		g, err := d.OrderPaidConsumerGroup()
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

func (d *diContainer) OrderPaidConsumerService() (service.OrderPaidConsumerService, error) {
	if d.orderPaidConsumerSvc == nil {
		c, err := d.OrderPaidConsumer()
		if err != nil {
			return nil, err
		}
		t, err := d.TelegramService()
		if err != nil {
			return nil, err
		}
		d.orderPaidConsumerSvc = orderPaidConsumerSvc.NewService(c, t)
	}
	return d.orderPaidConsumerSvc, nil
}

func (d *diContainer) OrderAssembledConsumerGroup() (sarama.ConsumerGroup, error) {
	if d.orderAssembledConsumerGroup == nil {
		cfg := config.AppConfig()
		g, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers(), cfg.OrderAssembledConsumer.GroupID(), cfg.OrderAssembledConsumer.Config())
		if err != nil {
			return nil, err
		}
		closer.AddNamed("Kafka consumer group (order assembled)", func(context.Context) error { return g.Close() })
		d.orderAssembledConsumerGroup = g
	}
	return d.orderAssembledConsumerGroup, nil
}

func (d *diContainer) OrderAssembledConsumer() (kafka.Consumer, error) {
	if d.orderAssembledConsumer == nil {
		g, err := d.OrderAssembledConsumerGroup()
		if err != nil {
			return nil, err
		}
		d.orderAssembledConsumer = kafkaConsumer.NewConsumer(
			g,
			[]string{config.AppConfig().OrderAssembledConsumer.Topic()},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.orderAssembledConsumer, nil
}

func (d *diContainer) OrderAssembledConsumerService() (service.OrderAssembledConsumerService, error) {
	if d.orderAssembledConsumerSvc == nil {
		c, err := d.OrderAssembledConsumer()
		if err != nil {
			return nil, err
		}
		t, err := d.TelegramService()
		if err != nil {
			return nil, err
		}
		d.orderAssembledConsumerSvc = orderAssembledConsumerSvc.NewService(c, t)
	}
	return d.orderAssembledConsumerSvc, nil
}
