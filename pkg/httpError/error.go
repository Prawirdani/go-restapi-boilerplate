package httpError

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorMessage interface {
	string | interface{}
}

type HttpError struct {
	Code    int
	Message ErrorMessage
}

func (e *HttpError) Error() string {
	return fmt.Sprint(e.Message)
}

func WrapError(err error) *HttpError {
	switch e := err.(type) {
	case *json.UnmarshalTypeError:
		return BadRequest(e.Error())
	default:
		if definedErr, ok := err.(*HttpError); ok {
			return definedErr
		}
		return InternalServer(err.Error())
	}
}

func BuildHttpError(code int) func(ErrorMessage) *HttpError {
	return func(msg ErrorMessage) *HttpError {
		return &HttpError{Code: code, Message: msg}
	}
}

var (
	InternalServer    = BuildHttpError(http.StatusInternalServerError)
	BadRequest        = BuildHttpError(http.StatusBadRequest)
	Unauthorized      = BuildHttpError(http.StatusUnauthorized)
	NotFound          = BuildHttpError(http.StatusNotFound)
	MethodNotAllowed  = BuildHttpError(http.StatusMethodNotAllowed)
	RequestValidation = parseValidationError
)

func parseValidationError(err error) *HttpError {
	switch e := err.(type) {
	case validator.ValidationErrors:
		errors := ValidatorError(e)
		return BadRequest(errors)
	default:
		return BadRequest(err)
	}
}

// it's for go-playground/validator/v10
type ValidationErrMsg map[string]interface{}

func ValidatorError(err validator.ValidationErrors) ValidationErrMsg {
	errors := make(ValidationErrMsg)
	for _, e := range err {
		field := strings.ToLower(e.Field())
		switch e.Tag() {
		case "required":
			errors[field] = fmt.Sprintf("%s field is required", field)
		case "eqfield":
			if field == "repeat_password" {
				errors[field] = "password don't match"
			}
		default:
			errors[field] = e.Error()
		}
	}
	return errors
}
