package health

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/gostock/internal/core/ports"
	"go.uber.org/fx"
)

type Handler struct {
	service ports.HealthService
	router  fiber.Router
}

type HandlerParams struct {
	fx.In
	Router  fiber.Router
	Service ports.HealthService
}

func New(params HandlerParams) (*Handler, error) {
	handler := Handler{params.Service, params.Router}

	handler.router.Get("/health", handler.Health)

	return &handler, nil
}

func (handler *Handler) Health(c fiber.Ctx) error {
	response := handler.service.Check(c.Context())
	status := fiber.StatusOK

	if response.Status != "UP" {
		status = fiber.StatusServiceUnavailable
	}

	return c.Status(status).JSON(response)
}
