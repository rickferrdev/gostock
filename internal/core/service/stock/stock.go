package stock

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/rickferrdev/gostock/internal/core/domain"
	"github.com/rickferrdev/gostock/internal/core/ports"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Service struct {
	stock ports.StockRepository
}

type ServiceParams struct {
	fx.In
	Stock ports.StockRepository
}

func New(params ServiceParams) (ports.StockService, error) {
	return &Service{params.Stock}, nil
}

func (service *Service) ValidateStock(stock *domain.Stock) error {
	if stock.Name == "" {
		return ports.NewBadRequestError(nil)
	}

	if stock.Capacity == 0 {
		return ports.NewBadRequestError(nil)
	}

	if err := service.ValidateID(stock.ID); err != nil {
		return ports.NewBadRequestError(err)
	}

	return nil
}

func (service *Service) ValidateID(id string) error {
	if err := uuid.Validate(id); err != nil {
		return ports.NewBadRequestError(err)
	}

	return nil
}

func (service *Service) Errors(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ports.NewNotFoundError(err)
	default:
		return ports.NewInternalError(err)
	}
}

func (service *Service) All(ctx context.Context) ([]*domain.Stock, error) {
	list, err := service.stock.All(ctx)
	if err != nil {
		return nil, service.Errors(err)
	}

	if len(list) == 0 {
		return nil, ports.NewNotFoundError(nil)
	}

	return list, nil
}

func (service *Service) ByID(ctx context.Context, id string) (*domain.Stock, error) {
	if err := service.ValidateID(id); err != nil {
		return nil, err
	}

	stock, err := service.stock.ByID(ctx, id)
	if err != nil {
		return nil, service.Errors(err)
	}

	return stock, nil
}

func (service *Service) Create(ctx context.Context, stock *domain.Stock) (*domain.Stock, error) {
	if stock.Name == "" {
		return nil, ports.NewBadRequestError(nil)
	}

	if stock.Capacity == 0 {
		return nil, ports.NewBadRequestError(nil)
	}

	empty := domain.Stock{
		ID:       uuid.NewString(),
		Name:     stock.Name,
		Capacity: stock.Capacity,
		Items:    make([]*domain.Product, 0),
	}

	if err := service.stock.Create(ctx, &empty); err != nil {
		return nil, service.Errors(err)
	}

	return &empty, nil
}

func (service *Service) Occupancy(ctx context.Context, id string) (uint64, error) {
	if err := service.ValidateID(id); err != nil {
		return 0, err
	}

	stock, err := service.ByID(ctx, id)
	if err != nil {
		return 0, service.Errors(err)
	}

	return stock.Occupancy(), nil
}

func (service *Service) AvailableSpace(ctx context.Context, id string) (uint64, error) {
	if err := service.ValidateID(id); err != nil {
		return 0, err
	}

	stock, err := service.ByID(ctx, id)
	if err != nil {
		return 0, service.Errors(err)
	}

	return stock.AvailableSpace(), nil
}

func (service *Service) Update(ctx context.Context, stock *domain.Stock) error {
	if err := service.ValidateID(stock.ID); err != nil {
		return err
	}

	if _, err := service.stock.ByID(ctx, stock.ID); err != nil {
		return service.Errors(err)
	}

	err := service.stock.Update(ctx, stock)
	if err != nil {
		return service.Errors(err)
	}

	return nil
}

func (service *Service) Remove(ctx context.Context, id string) error {
	if err := service.ValidateID(id); err != nil {
		return err
	}

	_, err := service.stock.ByID(ctx, id)
	if err != nil {
		return service.Errors(err)
	}

	if err := service.stock.Delete(ctx, id); err != nil {
		return service.Errors(err)
	}

	return nil
}
