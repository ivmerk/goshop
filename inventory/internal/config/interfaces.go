package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryGRPCConfig interface {
	Address() string
}

type IAMGRPCConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}
