package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

type structValidator struct {
	validator *validator.Validate
}

var Provide = fx.Provide(New)

func New() fiber.StructValidator {
	return &structValidator{validator: validator.New()}
}

func (val *structValidator) Validate(out any) error {
	return val.validator.Struct(out)
}
