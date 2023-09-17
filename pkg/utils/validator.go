package utils

import (
	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func ValidateRequest(request interface{}) error {
	return v.Struct(request)
	// slog.Error("utils.validator", "error", err)
}
