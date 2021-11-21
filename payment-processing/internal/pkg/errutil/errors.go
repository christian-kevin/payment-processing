package errutil

import "errors"

var (
	ErrServerError               = errors.New("server error")
	ErrInvalidParam              = errors.New("invalid param")
	ErrContextValueNotFound      = errors.New("value is invalid")
)
