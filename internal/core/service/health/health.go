package health

import (
	"context"

	"github.com/rickferrdev/gostock/internal/core/ports"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

type Service struct {
	sqlDB *bun.DB
}

var Provide = fx.Provide(New)

type ServiceParams struct {
	fx.In
	SqlDB *bun.DB
}

func New(params ServiceParams) (ports.HealthService, error) {
	service := Service{sqlDB: params.SqlDB}

	return &service, nil
}

func (service *Service) Check(ctx context.Context) *ports.HealthResponse {
	response := ports.HealthResponse{
		Status:   "UP",
		Services: map[string]any{},
	}

	if err := service.sqlDB.PingContext(ctx); err != nil {
		response.Services["database"] = "DOWN"
		response.Status = "DEGRADED"
	} else {
		response.Services["database"] = "UP"
	}

	return &response
}
