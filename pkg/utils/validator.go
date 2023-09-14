package utils

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func ValidateRequest(request interface{}) error {
	err := v.Struct(request)
	if err != nil {
		slog.Error("utils.validator", "error", err)
		return err
	}
	return nil
}
