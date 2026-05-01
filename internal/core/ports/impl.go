package ports

import (
	"context"

	"github.com/rickferrdev/gostock/internal/core/domain"
)

//go:generate mockgen -source=$GOFILE -destination=../../tests/mocks/ports.go -package=mocks
type StockRepository interface {
	All(ctx context.Context) ([]*domain.Stock, error)
	ByID(ctx context.Context, id string) (*domain.Stock, error)

	CreateAtomic(ctx context.Context, stock *domain.Stock) error

	UpdateAtomic(ctx context.Context, stock *domain.Stock) error
	DeleteAtomic(ctx context.Context, id string) error
}

type StockService interface {
	ValidateStock(stock *domain.Stock) error

	All(ctx context.Context) ([]*domain.Stock, error)
	ByID(ctx context.Context, id string) (*domain.Stock, error)

	Create(ctx context.Context, stock *domain.Stock) (*domain.Stock, error)

	Occupancy(ctx context.Context, id string) (uint64, error)
	AvailableSpace(ctx context.Context, id string) (uint64, error)

	Update(ctx context.Context, stock *domain.Stock) error
	Remove(ctx context.Context, id string) error
}

type ProductRepository interface {
	ByStockID(ctx context.Context, id string) ([]*domain.Product, error)
	ByID(ctx context.Context, id string) (*domain.Product, error)

	CreateAtomic(ctx context.Context, product *domain.Product) error

	UpdateAtomic(ctx context.Context, product *domain.Product) error
	DeleteAtomic(ctx context.Context, id string) error
}

type ProductService interface {
	ValidateProduct(product *domain.Product) error
	ValidateID(id string) error
	Errors(err error) error

	ByStockID(ctx context.Context, id string) ([]*domain.Product, error)
	ByID(ctx context.Context, id string) (*domain.Product, error)

	Create(ctx context.Context, product *domain.Product) (*domain.Product, error)

	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id string) error
}

type HealthService interface {
	Check(ctx context.Context) *HealthResponse
}
