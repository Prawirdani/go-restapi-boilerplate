package utils

import (
	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func ValidateRequest(request interface{}) error {
	err := v.Struct(request)
	if err != nil {
		return err
	}
	return nil
}
