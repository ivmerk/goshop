package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/ivmerk/goshop/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger LoggerConfig
	GRPC   InventoryGRPCConfig
	Mongo  MongoConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	logger, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	grpc, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	mongo, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger: logger,
		GRPC:   grpc,
		Mongo:  mongo,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
