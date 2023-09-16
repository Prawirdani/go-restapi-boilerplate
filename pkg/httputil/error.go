package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

type errMessage interface {
	string | interface{}
}

type apiError struct {
	Code    int
	Message errMessage
}

func (e *apiError) Error() string {
	return fmt.Sprint(e.Message)
}

func buildError(code int) func(errMessage) *apiError {
	return func(msg errMessage) *apiError {
		return &apiError{Code: code, Message: msg}
	}
}

var (
	ErrInternalServer   = buildError(http.StatusInternalServerError)
	ErrBadRequest       = buildError(http.StatusBadRequest)
	ErrUnauthorized     = buildError(http.StatusUnauthorized)
	ErrNotFound         = buildError(http.StatusNotFound)
	ErrMethodNotAllowed = buildError(http.StatusMethodNotAllowed)
)


func parseErrors(err error) *apiError {
	// By Error Type
	switch e := err.(type) {
	case validator.ValidationErrors:
		return parseValidationError(e)
	case *json.UnmarshalTypeError:
		return parseJsonError(e)
	case *mysql.MySQLError:
		return parseMysqlError(e)
	default:
		if httpError, ok := err.(*apiError); ok {
			return httpError
		}
		return ErrInternalServer(err.Error())
	}
}

// For go-playground/validator/v10 package
func parseValidationError(err validator.ValidationErrors) *apiError {
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
		default:
			errors[field] = errField.Error()
		}
	}
	return ErrBadRequest(errors)
}

func parseJsonError(err *json.UnmarshalTypeError) *apiError {
	return ErrBadRequest(err.Error())
}

func parseMysqlError(err *mysql.MySQLError) *apiError {
	if err.Number == 1062 {
		// Define your models duplicate entry key here
		return ErrBadRequest("{Key} already exists")
	}
	return ErrInternalServer(err)
}
