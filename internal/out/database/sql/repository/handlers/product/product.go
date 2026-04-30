package product

import (
	"context"
	"database/sql"

	"github.com/rickferrdev/gostock/internal/core/domain"
	"github.com/rickferrdev/gostock/internal/core/ports"
	"github.com/rickferrdev/gostock/internal/out/database/sql/repository/schema"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Repository struct {
	db *bun.DB
}

func New(db *bun.DB) (ports.ProductRepository, error) {
	return &Repository{db}, nil
}

func (repository *Repository) ByStockID(ctx context.Context, id string) ([]*domain.Product, error) {
	var models []*schema.Product

	if err := repository.db.NewSelect().
		Model(&models).
		Where("stock_id = ?", id).
		Scan(ctx); err != nil {
		return nil, err
	}

	domain := make([]*domain.Product, 0, len(models))

	for _, model := range models {
		convertion, err := model.ToDomain()
		if err != nil {
			return nil, err
		}

		domain = append(domain, convertion)
	}

	return domain, nil

}

func (repository *Repository) ByID(ctx context.Context, id string) (*domain.Product, error) {
	var model schema.Product

	if err := repository.db.NewSelect().
		Model(&model).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		return nil, err
	}

	return model.ToDomain()
}

func (repository *Repository) CreateAtomic(ctx context.Context, product *domain.Product) error {
	return repository.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var stock schema.Stock
		if err := tx.NewSelect().
			Model(&stock).
			Column("capacity").
			Where("id = ?", product.StockID).
			Scan(ctx); err != nil {
			return err
		}

		var occupancy uint64
		if err := tx.NewSelect().
			Model((*schema.Product)(nil)).
			ColumnExpr("COALESCE(SUM(qtd), 0)").
			Where("stock_id = ?", product.StockID).
			Scan(ctx, &occupancy); err != nil {
			return err
		}
		if occupancy+product.Qtd > stock.Capacity {
			return ports.NewCapacityExceeded(nil)
		}

		model, err := schema.FromProduct(product)
		if err != nil {
			return err
		}

		_, err = tx.NewInsert().Model(model).Exec(ctx)
		return err
	})
}

func (repository *Repository) UpdateAtomic(ctx context.Context, product *domain.Product) error {
	return repository.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var stock schema.Stock
		if err := tx.NewSelect().
			Model(&stock).
			Column("capacity").
			Where("id = ?", product.StockID).
			Scan(ctx); err != nil {
			return err
		}

		var old schema.Product
		if err := tx.NewSelect().
			Model(&old).
			Where("id = ?", product.ID).
			Scan(ctx); err != nil {
			return err
		}

		var occupancy uint64
		if err := tx.NewSelect().
			Model((*schema.Product)(nil)).
			ColumnExpr("COALESCE(SUM(qtd), 0)").
			Where("stock_id = ?", product.StockID).
			Scan(ctx, &occupancy); err != nil {
			return err
		}

		var total uint64

		if product.StockID == old.StockID {
			total = (occupancy - old.Qtd) + product.Qtd
		} else {
			total = occupancy + product.Qtd
		}

		if total > stock.Capacity {
			return ports.NewCapacityExceeded(nil)
		}

		model, err := schema.FromProduct(product)
		if err != nil {
			return err
		}

		if _, err = tx.NewUpdate().Model(model).WherePK().Exec(ctx); err != nil {
			return err
		}
		return err
	})
}

func (repository *Repository) Delete(ctx context.Context, id string) error {
	return repository.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		res, err := repository.db.NewDelete().
			Model((*schema.Product)(nil)).
			Where("id = ?", id).
			Exec(ctx)

		if err != nil {
			return err
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return sql.ErrNoRows
		}

		return nil
	})
}
