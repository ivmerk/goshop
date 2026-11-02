package config

import "time"

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
