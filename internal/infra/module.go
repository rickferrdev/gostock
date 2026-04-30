package infra

import (
	"github.com/rickferrdev/gostock/internal/infra/log"
	"github.com/rickferrdev/gostock/internal/infra/server"
	"github.com/rickferrdev/gostock/internal/infra/sql"
	"go.uber.org/fx"
)

var Module = fx.Module("infra", log.Provide, server.Provide, sql.Provide)
