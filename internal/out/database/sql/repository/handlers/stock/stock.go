package stock

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

func New(db *bun.DB) (ports.StockRepository, error) {
	return &Repository{db}, nil
}

func (repository *Repository) All(ctx context.Context) ([]*domain.Stock, error) {
	var models []*schema.Stock

	if err := repository.db.NewSelect().
		Model(&models).
		Relation("Items").
		Scan(ctx); err != nil {
		return nil, err
	}

	domain := make([]*domain.Stock, 0, len(models))

	for _, model := range models {
		conversion, err := model.ToDomain()
		if err != nil {
			return nil, err
		}

		domain = append(domain, conversion)
	}

	return domain, nil

}

func (repository *Repository) ByID(ctx context.Context, id string) (*domain.Stock, error) {
	var model schema.Stock

	if err := repository.db.NewSelect().
		Model(&model).
		Relation("Items").
		Where("id = ?", id).
		Scan(ctx); err != nil {
		return nil, err
	}

	return model.ToDomain()
}

func (repository *Repository) CreateAtomic(ctx context.Context, stock *domain.Stock) error {
	return repository.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		model, err := schema.FromStock(stock)
		if err != nil {
			return err
		}

		if _, err = tx.NewInsert().Model(model).Exec(ctx); err != nil {
			return err
		}

		if len(model.Items) > 0 {
			_, err = tx.NewInsert().Model(&model.Items).Exec(ctx)
		}

		return err
	})
}

func (repository *Repository) UpdateAtomic(ctx context.Context, stock *domain.Stock) error {
	return repository.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		model, err := schema.FromStock(stock)
		if err != nil {
			return err
		}

		result, err := tx.NewUpdate().Model(model).WherePK().Exec(ctx)
		if err != nil {
			return err
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return sql.ErrNoRows
		}

		return nil
	})
}

func (repository *Repository) DeleteAtomic(ctx context.Context, id string) error {
	return repository.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewDelete().Model((*schema.Stock)(nil)).Where("id = ?", id).Exec(ctx)
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
