package env

import (
	"net"
	"net/url"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host          string `env:"POSTGRES_HOST,required"`
	Port          string `env:"POSTGRES_PORT,required"`
	Database      string `env:"POSTGRES_DB,required"`
	User          string `env:"POSTGRES_USER,required"`
	Password      string `env:"POSTGRES_PASSWORD,required"`
	SSLMode       string `env:"POSTGRES_SSL_MODE,required"`
	MigrationsDir string `env:"MIGRATION_DIRECTORY,required"`
}

type postgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresConfig() (*postgresConfig, error) {
	var raw postgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &postgresConfig{raw: raw}, nil
}

func (cfg *postgresConfig) URI() string {
	u := url.URL{
		Scheme:   "postgres",
		Host:     net.JoinHostPort(cfg.raw.Host, cfg.raw.Port),
		Path:     "/" + cfg.raw.Database,
		User:     url.UserPassword(cfg.raw.User, cfg.raw.Password),
		RawQuery: "sslmode=" + url.QueryEscape(cfg.raw.SSLMode),
	}
	return u.String()
}

func (cfg *postgresConfig) MigrationsDir() string {
	return cfg.raw.MigrationsDir
}
