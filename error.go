package gin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	ErrBadRequest          = func() *Error { return NewError(http.StatusBadRequest) }
	ErrUnauthorized        = func() *Error { return NewError(http.StatusUnauthorized) }
	ErrForbidden           = func() *Error { return NewError(http.StatusForbidden) }
	ErrNotFound            = func() *Error { return NewError(http.StatusNotFound) }
	ErrConflict            = func() *Error { return NewError(http.StatusConflict) }
	ErrUnprocessableEntity = func() *Error { return NewError(http.StatusUnprocessableEntity) }
	ErrInternalServerError = func() *Error { return NewError(http.StatusInternalServerError) }

	ErrorTransform = DefaultErrorTransform
)

type Error struct {
	Code     int
	Message  string
	Data     interface{}
	Internal error
}

type FieldError struct {
	Namespace string `json:"namespace,omitempty"`
	Field     string `json:"field,omitempty"`
	Tag       string `json:"tag,omitempty"`
	Value     string `json:"value,omitempty"`
	Message   string `json:"message,omitempty"`
}

type ErrorTransformFunc func(err error) *Error

func NewError(code int) *Error {
	return &Error{
		Code:    code,
		Message: http.StatusText(code),
	}
}

func (e *Error) Error() string {
	if e.Internal == nil {
		return fmt.Sprintf("code=%d, message=%v, data=%v", e.Code, e.Message, e.Data)
	}
	return fmt.Sprintf("code=%d, message=%v, data=%v, internal=%v", e.Code, e.Message, e.Data, e.Internal)
}

func (e *Error) Unwrap() error {
	return e.Internal
}

func (e *Error) SetMessage(message string) *Error {
	e.Message = message
	return e
}

func (e *Error) SetData(data interface{}) *Error {
	e.Data = data
	return e
}

func (e *Error) SetInternal(err error) *Error {
	e.Internal = err
	return e
}

func (e *Error) AddInternal(err error) *Error {
	if e.Internal == nil {
		e.Internal = err
	} else if err != nil {
		e.Internal = fmt.Errorf("%w; %w", e.Internal, err)
	}
	return e
}

func (e *Error) MarshalJSON() ([]byte, error) {
	err := struct {
		Code      int         `json:"code,omitempty"`
		Message   string      `json:"message,omitempty"`
		Data      interface{} `json:"data,omitempty"`
		Developer string      `json:"developer_message,omitempty"`
	}{
		Code:    e.Code,
		Message: e.Message,
	}

	if gin.IsDebugging() && e.Internal != nil {
		err.Developer = e.Internal.Error()
	}

	return json.Marshal(err)
}

func DefaultErrorTransform(err error) (e *Error) {
	if errors.As(err, &e) {
		return
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		data := make([]FieldError, len(ve))
		for i, field := range ve {
			data[i] = FieldError{
				Namespace: field.StructNamespace(),
				Field:     field.Field(),
				Tag:       field.Tag(),
				Value:     field.Param(),
				Message:   field.Error(),
			}
		}
		e = ErrUnprocessableEntity().SetData(data)
	} else {
		e = ErrInternalServerError()
	}

	return
}
