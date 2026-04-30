package platform

import (
	"github.com/rickferrdev/gostock/internal/platform/validator"
	"go.uber.org/fx"
)

var Module = fx.Module("platform", validator.Provide)
