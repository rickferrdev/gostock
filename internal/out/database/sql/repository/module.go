package repository

import (
	"github.com/rickferrdev/gostock/internal/out/database/sql/repository/handlers/product"
	"github.com/rickferrdev/gostock/internal/out/database/sql/repository/handlers/stock"
	"go.uber.org/fx"
)

var Module = fx.Module("repository", product.Provide, stock.Provide)
