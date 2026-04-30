package schema

import (
	"github.com/google/uuid"
	"github.com/rickferrdev/gostock/internal/core/domain"
	"github.com/uptrace/bun"
)

type Stock struct {
	bun.BaseModel

	ID       string `bun:"type:uuid,pk"`
	Name     string `bun:"name"`
	Capacity uint64 `bun:"capacity"`

	Items []*Product `bun:"rel:has-many,join:id=stock_id"`
}

type Product struct {
	bun.BaseModel

	ID    string `bun:"type:uuid,pk"`
	Name  string `bun:"name,notnull"`
	Qtd   uint64 `bun:"qtd,notnull"`
	Price int64  `bun:"price,notnull"`

	StockID string `bun:"type:uuid,notnull"`

	Stock *Stock `bun:"rel:belongs-to,join:stock_id=id"`
}

func (stock *Stock) ToDomain() (*domain.Stock, error) {
	if err := uuid.Validate(stock.ID); err != nil {
		return nil, err
	}

	items := make([]*domain.Product, 0, len(stock.Items))

	for _, i := range stock.Items {
		product, err := i.ToDomain()
		if err != nil {
			return nil, err
		}

		items = append(items, product)
	}

	return &domain.Stock{
		ID:       stock.ID,
		Name:     stock.Name,
		Capacity: stock.Capacity,
		Items:    items,
	}, nil
}

func (product *Product) ToDomain() (*domain.Product, error) {
	if err := uuid.Validate(product.ID); err != nil {
		return nil, err
	}

	if err := uuid.Validate(product.StockID); err != nil {
		return nil, err
	}

	return &domain.Product{
		ID:      product.ID,
		Name:    product.Name,
		Qtd:     product.Qtd,
		Price:   product.Price,
		StockID: product.StockID,
	}, nil
}

func FromStock(stock *domain.Stock) (*Stock, error) {
	if err := uuid.Validate(stock.ID); err != nil {
		return nil, err
	}

	items := make([]*Product, 0, len(stock.Items))

	for _, i := range stock.Items {
		product, err := FromProduct(i)
		if err != nil {
			return nil, err
		}

		items = append(items, product)
	}

	return &Stock{
		ID:       stock.ID,
		Name:     stock.Name,
		Capacity: stock.Capacity,
		Items:    items,
	}, nil
}

func FromProduct(product *domain.Product) (*Product, error) {
	if err := uuid.Validate(product.ID); err != nil {
		return nil, err
	}

	if err := uuid.Validate(product.StockID); err != nil {
		return nil, err
	}

	return &Product{
		ID:      product.ID,
		Name:    product.Name,
		Qtd:     product.Qtd,
		Price:   product.Price,
		StockID: product.StockID,
	}, nil
}

// func Invoke(db *bun.DB) error {
// 	ctx := context.Background()

// 	_, err := db.NewCreateTable().
// 		Model((*Stock)(nil)).
// 		IfNotExists().
// 		Exec(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = db.NewCreateTable().
// 		Model((*Product)(nil)).
// 		WithForeignKeys().
// 		IfNotExists().
// 		Exec(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
