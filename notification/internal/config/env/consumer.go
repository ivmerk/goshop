package env

import (
	"github.com/IBM/sarama"
)

// consumerConfig — общая реализация настроек Kafka consumer'а для разных топиков.
type consumerConfig struct {
	topic   string
	groupID string
}

func (cfg *consumerConfig) Topic() string {
	return cfg.topic
}

func (cfg *consumerConfig) GroupID() string {
	return cfg.groupID
}

func (cfg *consumerConfig) Config() *sarama.Config {
	c := sarama.NewConfig()
	c.Version = sarama.V4_0_0_0
	c.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	c.Consumer.Offsets.Initial = sarama.OffsetOldest

	return c
}
