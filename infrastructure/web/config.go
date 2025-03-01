package web

import (
	"strings"

	"github.com/goda6565/nexus-user-auth/pkg/utils"
)

type Config struct {
	Host             string
	Port             string
	CorsAllowOrigins []string
}

func NewConfigWeb() *Config {
	return &Config{
		Host:             utils.GetEnvDefault("WEB_HOST", "localhost"),
		Port:             utils.GetEnvDefault("WEB_PORT", "8080"),
		CorsAllowOrigins: strings.Split(utils.GetEnvDefault("CORS_ALLOW_ORIGINS", "http://localhost:3000"), ","),
	}
}
