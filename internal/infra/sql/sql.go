package sql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rickferrdev/gostock/internal/config/env"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

func New(life fx.Lifecycle, env *env.Env) (*bun.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dir := filepath.Dir(env.APP_DATABASE_URI)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	fl := fmt.Sprintf("file:%v?cache=shared&mode=rwc", env.APP_DATABASE_URI)

	sqlite, err := sql.Open(sqliteshim.ShimName, fl)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqlite, sqlitedialect.New())

	if _, err = db.ExecContext(ctx, "PRAGMA journal_mode=WAL;"); err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	life.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				return db.Close()
			},
		},
	)

	return db, nil
}
