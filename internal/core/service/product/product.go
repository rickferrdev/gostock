package product

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
	product ports.ProductRepository
	stock   ports.StockRepository
}

type ServiceParams struct {
	fx.In
	Product ports.ProductRepository
	Stock   ports.StockRepository
}

func New(params ServiceParams) (ports.ProductService, error) {
	return &Service{params.Product, params.Stock}, nil
}

func (service *Service) ValidateProduct(product *domain.Product) error {
	if product.Name == "" {
		return ports.NewBadRequestError(nil)
	}

	if err := uuid.Validate(product.ID); err != nil {
		return ports.NewBadRequestError(err)
	}

	if err := uuid.Validate(product.StockID); err != nil {
		return ports.NewBadRequestError(err)
	}

	if product.Price < 0 {
		return ports.NewBadRequestError(nil)
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

func (service *Service) ByStockID(ctx context.Context, id string) ([]*domain.Product, error) {
	if err := service.ValidateID(id); err != nil {
		return nil, err
	}

	list, err := service.product.ByStockID(ctx, id)
	if err != nil {
		return nil, service.Errors(err)
	}

	if len(list) == 0 {
		return nil, ports.NewNotFoundError(nil)
	}

	return list, nil
}

func (service *Service) ByID(ctx context.Context, id string) (*domain.Product, error) {
	if err := service.ValidateID(id); err != nil {
		return nil, err
	}

	product, err := service.product.ByID(ctx, id)
	if err != nil {
		return nil, service.Errors(err)
	}

	return product, nil
}

func (service *Service) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	if err := uuid.Validate(product.StockID); err != nil {
		return nil, ports.NewBadRequestError(err)
	}

	if product.ID != "" {
		return nil, ports.NewConflictError(nil)
	}

	if product.StockID == "" {
		return nil, ports.NewConflictError(nil)
	}

	pro := domain.Product{
		ID:      uuid.NewString(),
		Name:    product.Name,
		Qtd:     product.Qtd,
		Price:   product.Price,
		StockID: product.StockID,
	}

	if err := service.product.CreateAtomic(ctx, &pro); err != nil {
		return nil, service.Errors(err)
	}

	return &pro, nil
}

func (service *Service) Update(ctx context.Context, product *domain.Product) error {
	if err := service.ValidateProduct(product); err != nil {
		return err
	}

	oldProduct, err := service.product.ByID(ctx, product.ID)
	if err != nil {
		return service.Errors(err)
	}

	stock, err := service.stock.ByID(ctx, product.StockID)
	if err != nil {
		return service.Errors(err)
	}

	if oldProduct.StockID == product.StockID {
		realAvailableSpace := stock.AvailableSpace() + oldProduct.Qtd
		if realAvailableSpace < product.Qtd {
			return ports.NewCapacityExceeded(nil)
		}
	} else {
		if !stock.CanAdd(product.Qtd) {
			return ports.NewCapacityExceeded(nil)
		}
	}

	pro := domain.Product{
		ID:      product.ID,
		Name:    product.Name,
		Qtd:     product.Qtd,
		Price:   product.Price,
		StockID: product.StockID,
	}

	if err := service.product.UpdateAtomic(ctx, &pro); err != nil {
		return service.Errors(err)
	}

	return nil
}

func (service *Service) Delete(ctx context.Context, id string) error {
	if err := service.ValidateID(id); err != nil {
		return err
	}

	if err := service.product.Delete(ctx, id); err != nil {
		return service.Errors(err)
	}

	return nil
}
