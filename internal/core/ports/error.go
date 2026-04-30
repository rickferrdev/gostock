package ports

import (
	"fmt"
	"net/http"
)

type Code string
type Message string

type Error struct {
	Status  int
	Message Message
	Code    Code
	Err     error
}

const (
	CodeInternal        Code = "INTERNAL_SERVER_ERROR"
	CodeBadRequest      Code = "BAD_REQUEST"
	CodeNotFound        Code = "RESOURCE_NOT_FOUND"
	CodeUnauthenticated Code = "UNAUTHENTICATED"
	CodeForbidden       Code = "PERMISSION_DENIED"
	CodeConflict        Code = "RESOURCE_CONFLICT"

	CodeUserNotFound    Code = "USER_NOT_FOUND"
	CodeInvalidPassword Code = "INVALID_PASSWORD"
	CodeEmailRegistered Code = "EMAIL_ALREADY_REGISTERED"
	CodeTokenExpired    Code = "AUTH_TOKEN_EXPIRED"

	CodeDatabaseError Code = "DATABASE_CONNECTION_FAILED"
	CodeTimeout       Code = "OPERATION_TIMEOUT"

	MsgInternal          Message = "an unexpected internal error occurred"
	MsgNotFound          Message = "the requested resource was not found"
	MsgInvalidInput      Message = "the provided data is invalid"
	MsgUnauthorized      Message = "you must be authenticated to access this resource"
	MsgForbidden         Message = "you do not have permission to perform this action"
	MsgUserConflict      Message = "user with these details already exists"
	MsgBadRequest        Message = "invalid request parameters"
	MsgStockNotEmpty     Message = "cannot remove a stock that still has registered products"
	MsgCapacityExceeded  Message = "the resource's capacity was exceeded"
	MsgProcessingMessage Message = "error processing the standard submission message"
)

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s:%d] %s: %v", e.Code, e.Status, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s:%d] %s", e.Code, e.Status, e.Message)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func NewError(code Code, message Message, status int, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

func NewNotFoundError(err error) *Error {
	return NewError(CodeNotFound, MsgNotFound, http.StatusNotFound, err)
}

func NewBadRequestError(err error) *Error {
	return NewError(CodeBadRequest, MsgBadRequest, http.StatusBadRequest, err)
}

func NewInternalError(err error) *Error {
	return NewError(CodeInternal, MsgInternal, http.StatusInternalServerError, err)
}

func NewCapacityExceeded(err error) *Error {
	return NewError(CodeConflict, MsgCapacityExceeded, http.StatusConflict, err)
}

func NewConflictError(err error) *Error {
	return NewError(CodeConflict, MsgUserConflict, http.StatusConflict, err)
}

func NewStockNotEmptyError(err error) *Error {
	return NewError(CodeConflict, MsgStockNotEmpty, http.StatusConflict, err)
}
