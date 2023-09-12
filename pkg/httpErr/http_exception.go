package httpErr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	json_util "github.com/prawirdani/go-restapi-boilerplate/pkg/json"
)

func ExceptionIfErr(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *HttpExceptionError: // Our Custom Error
		json_util.Send(w, e.Code, e.Message)
	case *json.UnmarshalTypeError:
		json_util.Send(w, http.StatusBadRequest, e.Error())
	default:
		/* Add your logger here */
		json_util.Send(w, http.StatusInternalServerError, err.Error())
	}
}

/* Used for expected/trackable error such as no data found, invalid request body from client and invalid authorization */

type ErrorMessage interface {
	string | interface{}
}

type HttpExceptionError struct {
	Code    int
	Message ErrorMessage
}

func (e *HttpExceptionError) Error() string {
	return fmt.Sprint(e.Message)
}

func BuildHttpError(code int) func(ErrorMessage) *HttpExceptionError {
	return func(msg ErrorMessage) *HttpExceptionError {
		return &HttpExceptionError{Code: code, Message: msg}
	}
}

var (
	ErrInteralServer     = BuildHttpError(500)
	ErrBadRequest        = BuildHttpError(400)
	ErrUnauthorized      = BuildHttpError(401)
	ErrNotFound          = BuildHttpError(404)
	ErrRequestValidation = func(err error) *HttpExceptionError {
		switch e := err.(type) {
		case validator.ValidationErrors:
			errors := ValidatorError(e)
			return BuildHttpError(400)(errors)
		default:
			return ErrBadRequest(err)
		}
	}
)

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
