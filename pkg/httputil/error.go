package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

type errCause interface {
	string | interface{}
}

type Error struct {
	Code    int
	Message errCause
}

func (e *Error) Error() string {
	return fmt.Sprint(e.Message)
}

func buildError(code int) func(errCause) *Error {
	return func(msg errCause) *Error {
		return &Error{Code: code, Message: msg}
	}
}

var (
	ErrInternalServer   = buildError(http.StatusInternalServerError)
	ErrBadRequest       = buildError(http.StatusBadRequest)
	ErrUnauthorized     = buildError(http.StatusUnauthorized)
	ErrNotFound         = buildError(http.StatusNotFound)
	ErrMethodNotAllowed = buildError(http.StatusMethodNotAllowed)
)

func parseErrors(err error) *Error {
	// By Error String
	switch {
	case strings.Contains(err.Error(), "EOF"):
		return ErrBadRequest("Invalid JSON request body format")
	}

	// By Error Type
	switch e := err.(type) {
	case validator.ValidationErrors:
		return parseValidationError(e)
	case *json.UnmarshalTypeError:
		return parseJsonError(e)
	case *pgconn.PgError:
		return parsePostgreError(e)
	default:
		if httpError, ok := err.(*Error); ok {
			return httpError
		}
		return ErrInternalServer(err.Error())
	}
}

// For go-playground/validator/v10 package
func parseValidationError(err validator.ValidationErrors) *Error {
	errors := make(map[string]interface{})
	for _, errField := range err {
		field := strings.ToLower(errField.Field())
		switch errField.Tag() {
		case "required":
			errors[field] = fmt.Sprintf("%s field is required", field)
		case "eqfield":
			if field == "repeat_password" {
				errors[field] = "password don't match"
			}
		case "email":
			errors[field] = "invalid email format"
		default:
			errors[field] = errField.Error()
		}
	}
	return ErrBadRequest(errors)
}

func parseJsonError(err *json.UnmarshalTypeError) *Error {
	if strings.Contains(err.Error(), "unmarshal") {
		return ErrBadRequest(fmt.Sprintf("Type mismatch at %s, Expected type %s, Got %s", err.Field, err.Type, err.Value))
	}
	return ErrBadRequest(err.Error())
}

func parsePostgreError(err *pgconn.PgError) *Error {
	if err.Code == "23505" { // Duplicate Key Error Code
		return ErrBadRequest(err.Detail)
	}
	return ErrInternalServer(err)
}
