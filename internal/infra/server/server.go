package server

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/favicon"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/rickferrdev/gostock/assets"
	swagger "github.com/rickferrdev/gostock/docs/swagger"
	"github.com/rickferrdev/gostock/internal/config/env"
	"github.com/rickferrdev/gostock/internal/core/ports"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type ServerParams struct {
	fx.In
	Validator fiber.StructValidator
	Env       *env.Env
}

func New(params ServerParams) (*fiber.App, fiber.Router) {
	app := fiber.New(fiber.Config{
		CaseSensitive:   true,
		ErrorHandler:    ErrorHandler,
		StructValidator: params.Validator,
	})

	if params.Env == nil {
		params.Env = &env.Env{
			APP_SERVER_ORIGIN_FRONT: "*",
		}
	}

	registerMiddlewares(app, params.Env)

	return app, app.Group("/api/v1")
}

func registerMiddlewares(app *fiber.App, env *env.Env) {
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Expiration:             30 * time.Second,
		SkipSuccessfulRequests: true,
		Max:                    10,
	}))
	app.Use(favicon.New(favicon.Config{
		FileSystem: assets.Favicon,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{env.APP_SERVER_ORIGIN_FRONT},
		AllowMethods: []string{fiber.MethodGet, fiber.MethodPost, fiber.MethodDelete, fiber.MethodPut, fiber.MethodPatch, fiber.MethodConnect},
	}))
	app.Use(swaggerui.New(swaggerui.Config{
		FileContent: swagger.SwaggerJSON,
		BasePath:    "/api/v1",
		Path:        "docs",
	}))
}

func ErrorHandler(c fiber.Ctx, err error) error {
	code := ports.CodeInternal

	status := fiber.StatusInternalServerError
	attr := []any{slog.String("type", "unknown")}

	var f *ports.Error
	var e *fiber.Error
	var message string

	switch {
	case errors.As(err, &f):
		status = f.Status
		code = f.Code
		message = string(f.Message)
		attr = []any{
			slog.String("type", "domain"),
			slog.String("code", string(f.Code)),
		}

	case errors.As(err, &e):
		status = e.Code
		message = e.Message
		attr = []any{
			slog.String("type", "fiber"),
			slog.Int("status", e.Code),
		}
	}

	slog.Error(
		fmt.Sprintf("[%s:%d] %s", code, status, message),
		slog.String("request_id", requestid.FromContext(c)),
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.Int("status", status),
		slog.Any("error", err),
		slog.Group("context", attr...),
	)

	return c.Status(status).JSON(fiber.Map{
		"code":    code,
		"message": message,
	})
}
