package env

import (
	"github.com/caarlos0/env/v11"
)

type orderPaidConsumerEnvConfig struct {
	TopicName string `env:"ORDER_PAID_TOPIC_NAME,required"`
	GroupID   string `env:"ORDER_PAID_CONSUMER_GROUP_ID,required"`
}

func NewOrderPaidConsumerConfig() (*consumerConfig, error) {
	var raw orderPaidConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &consumerConfig{topic: raw.TopicName, groupID: raw.GroupID}, nil
}
