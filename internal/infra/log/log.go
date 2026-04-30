package log

import (
	"fmt"
	"log/slog"

	"github.com/rickferrdev/gostock/internal/config/env"
	slogbetterstack "github.com/samber/slog-betterstack"

	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

func New(env *env.Env) *slog.Logger {
	logger := slog.New(
		slogbetterstack.Option{
			Token:    env.APP_BETTERSTACK_TOKEN,
			Endpoint: fmt.Sprintf("https://%s/", env.APP_BETTERSTACK_INGESTING_HOST),
		}.NewBetterstackHandler(),
	)

	slog.SetDefault(logger)
	return logger
}
