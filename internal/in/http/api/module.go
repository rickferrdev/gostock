package api

import (
	"github.com/rickferrdev/gostock/internal/in/http/api/handlers/health"
	"github.com/rickferrdev/gostock/internal/in/http/api/handlers/product"
	"github.com/rickferrdev/gostock/internal/in/http/api/handlers/stock"
	"go.uber.org/fx"
)

var Module = fx.Module("api", fx.Invoke(stock.New, product.New, health.New))
