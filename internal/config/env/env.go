package env

import (
	"github.com/rickferrdev/dotenv"
	"go.uber.org/fx"
)

type Env struct {
	APP_SERVER_PORT                string `env:"APP_SERVER_PORT"`
	APP_DATABASE_URI               string `env:"APP_DATABASE_URI"`
	APP_BETTERSTACK_TOKEN          string `env:"APP_BETTERSTACK_TOKEN"`
	APP_BETTERSTACK_INGESTING_HOST string `env:"APP_BETTERSTACK_INGESTING_HOST"`
	APP_SERVER_ORIGIN_FRONT        string `env:"APP_SERVER_ORIGIN_FRONT"`
}

var Provide = fx.Provide(New)

func New() (*Env, error) {
	var e Env

	dotenv.Collect()

	if err := dotenv.Unmarshal(&e); err != nil {
		return nil, err
	}

	return &e, nil
}
