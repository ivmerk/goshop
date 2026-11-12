package env

import (
	"github.com/caarlos0/env/v11"
)

type orderAssembledConsumerEnvConfig struct {
	TopicName string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
	GroupID   string `env:"ORDER_ASSEMBLED_CONSUMER_GROUP_ID,required"`
}

func NewOrderAssembledConsumerConfig() (*consumerConfig, error) {
	var raw orderAssembledConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &consumerConfig{topic: raw.TopicName, groupID: raw.GroupID}, nil
}
