package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderPaidProducerEnvConfig struct {
	TopicName string `env:"ORDER_PAID_TOPIC_NAME,required"`
}

type orderPaidProducerConfig struct {
	raw orderPaidProducerEnvConfig
}

func NewOrderPaidProducerConfig() (*orderPaidProducerConfig, error) {
	var raw orderPaidProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidProducerConfig{raw: raw}, nil
}

func (cfg *orderPaidProducerConfig) Topic() string {
	return cfg.raw.TopicName
}

func (cfg *orderPaidProducerConfig) Config() *sarama.Config {
	c := sarama.NewConfig()
	c.Version = sarama.V4_0_0_0
	c.Producer.Return.Successes = true
	c.Producer.RequiredAcks = sarama.WaitForAll // ждём подтверждения брокера

	return c
}
