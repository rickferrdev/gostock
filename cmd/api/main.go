// Package main is the entry point of the Product Manager API.
//
// @title           Product Manager API
// @version         1.0
// @description     REST API for managing stocks and their products.
//
// @contact.name    Rick Ferr
// @contact.url     https://github.com/rickferrdev
//
// @license.name    MIT
//
// @host            localhost:8080
// @BasePath        /api/v1
//
// @schemes         http https
// @servers.url     http://localhost:8080/api/v1
package main

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/gostock/internal/config"
	"github.com/rickferrdev/gostock/internal/config/env"
	"github.com/rickferrdev/gostock/internal/core/service"
	"github.com/rickferrdev/gostock/internal/in/http/api"
	"github.com/rickferrdev/gostock/internal/infra"
	"github.com/rickferrdev/gostock/internal/out/database/sql/repository"
	"github.com/rickferrdev/gostock/internal/platform"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		infra.Module,
		platform.Module,
		repository.Module,
		service.Module,
		api.Module,
		fx.Invoke(Start),
	).Run()
}

func Start(life fx.Lifecycle, app *fiber.App, env *env.Env, log *slog.Logger) {
	life.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				port := env.APP_SERVER_PORT
				if port == "" {
					port = "8080"
				}

				go func() {
					if err := app.Listen("0.0.0.0:" + port); err != nil {
						log.Error("error starting server", slog.String("error", err.Error()))
					}
				}()

				log.Info("server started successfully")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return app.ShutdownWithContext(ctx)
			},
		},
	)
}
