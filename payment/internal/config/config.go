package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/ivmerk/goshop/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger LoggerConfig
	GRPC   PaymentGRPCConfig
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

	grpc, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger: logger,
		GRPC:   grpc,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
