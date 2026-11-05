package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryGRPCConfig interface {
	Address() string
}

type OrderHTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type PaymentGRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	MigrationsDir() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}
