package service

import (
	"github.com/rickferrdev/gostock/internal/core/service/health"
	"github.com/rickferrdev/gostock/internal/core/service/product"
	"github.com/rickferrdev/gostock/internal/core/service/stock"
	"go.uber.org/fx"
)

var Module = fx.Module("service", stock.Provide, product.Provide, health.Provide)
