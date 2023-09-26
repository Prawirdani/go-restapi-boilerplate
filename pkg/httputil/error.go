package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
)

type RestError struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Cause  interface{} `json:"error"`
}

func (e *RestError) Error() string {
	return fmt.Sprint(e.Cause)
}

func buildRestError(statusCode int) func(interface{}) *RestError {
	return func(cause interface{}) *RestError {
		return &RestError{Code: statusCode, Status: http.StatusText(statusCode), Cause: cause}
	}
}

var (
	ErrNotFound         = buildRestError(http.StatusNotFound)
	ErrBadRequest       = buildRestError(http.StatusBadRequest)
	ErrUnauthorized     = buildRestError(http.StatusUnauthorized)
	ErrInternalServer   = buildRestError(http.StatusInternalServerError)
	ErrMethodNotAllowed = buildRestError(http.StatusMethodNotAllowed)
)

func parseErrors(err error) *RestError {
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
		return parsePostgresError(e)
	default:
		if httpError, ok := err.(*RestError); ok {
			return httpError
		}
		return ErrInternalServer(err.Error())
	}
}

// For go-playground/validator/v10 package
func parseValidationError(err validator.ValidationErrors) *RestError {
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

func parseJsonError(err *json.UnmarshalTypeError) *RestError {
	if strings.Contains(err.Error(), "unmarshal") {
		return ErrBadRequest(fmt.Sprintf("Type mismatch at %s, Expected type %s, Got %s", err.Field, err.Type, err.Value))
	}
	return ErrBadRequest(err.Error())
}

func parsePostgresError(err *pgconn.PgError) *RestError {
	if err.Code == "23505" { // Duplicate Key Error Code
		switch {
		case err.ConstraintName == "users_username_key":
			return ErrBadRequest("username already exist")
		case err.ConstraintName == "users_email_key":
			return ErrBadRequest("email already exist")
		default:
			return ErrBadRequest(err.Detail)
		}
	}
	return ErrInternalServer(err)
}
