package app

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServieName string `default:"webapp-tutorial-backend" envconfig:"SERVICE_NAME"`

	// infra
	DatabaseHost string `required:"true" envconfig:"DATABASE_HOST"`
	DatabasePort string `required:"true" envconfig:"DATABASE_PORT"`
	DatabaseUser string `required:"true" envconfig:"DATABASE_USER"`
	DatabasePass string `required:"true" envconfig:"DATABASE_PASS"`
	DatabaseName string `required:"true" envconfig:"DATABASE_NAME"`

	// User usecase
	UserTable string `required:"true" envconfig:"USER_TABLE"`

	// Auth usecase
	AuthSecret               string        `required:"true" envconfig:"AUTH_SECRET"`
	AuthAccessTokenDuration  time.Duration `default:"6h" envconfig:"AUTH_ACCESS_TOKEN_DURATION"`
	AuthRefreshTokenDuration time.Duration `default:"720h" envconfig:"AUTH_REFRESH_TOKEN_DURATION"`

	// Rest Presenter
	RestAuthSecureRefreshToken bool `required:"true" envconfig:"REST_AUTH_SECURE_REFRESH_TOKEN"`
}

func NewConfig() (*Config, error) {
	var env Config
	if err := envconfig.Process("", &env); err != nil {
		return nil, err
	}

	return &env, nil
}
