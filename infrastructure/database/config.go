package database

import (
	"github.com/goda6565/nexus-user-auth/pkg/utils"
)

type Config struct {
	Driver   string
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func NewConfigPostgres() *Config {
	return &Config{
		Driver:   utils.GetEnvDefault("DB_DRIVER", "postgres"),
		User:     utils.GetEnvDefault("DB_USER", "postgres"),
		Password: utils.GetEnvDefault("DB_PASSWORD", "postgres"),
		Host:     utils.GetEnvDefault("DB_HOST", "localhost"),
		Port:     utils.GetEnvDefault("DB_PORT", "5432"),
		Database: utils.GetEnvDefault("DB_NAME", "auth"),
	}
}

func NewConfigSQLite() *Config {
	return &Config{
		Database: utils.GetEnvDefault("DB_NAME", "auth.sqlite"),
	}
}