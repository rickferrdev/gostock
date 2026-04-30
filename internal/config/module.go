package config

import (
	"github.com/rickferrdev/gostock/internal/config/env"
	"go.uber.org/fx"
)

var Module = fx.Module("config", env.Provide)
