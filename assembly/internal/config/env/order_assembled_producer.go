package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderAssembledProducerEnvConfig struct {
	TopicName string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
}

type orderAssembledProducerConfig struct {
	raw orderAssembledProducerEnvConfig
}

func NewOrderAssembledProducerConfig() (*orderAssembledProducerConfig, error) {
	var raw orderAssembledProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderAssembledProducerConfig{raw: raw}, nil
}

func (cfg *orderAssembledProducerConfig) Topic() string {
	return cfg.raw.TopicName
}

func (cfg *orderAssembledProducerConfig) Config() *sarama.Config {
	c := sarama.NewConfig()
	c.Version = sarama.V4_0_0_0
	c.Producer.Return.Successes = true
	c.Producer.RequiredAcks = sarama.WaitForAll // ждём подтверждения брокера

	return c
}
