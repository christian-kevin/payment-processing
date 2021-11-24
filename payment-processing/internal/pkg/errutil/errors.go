package errutil

import "errors"

var (
	ErrServerError          = errors.New("server error")
	ErrInvalidParam         = errors.New("invalid param")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrContextValueNotFound = errors.New("value is invalid")
	ErrDuplicateRequest     = errors.New("duplicate request")
	ErrWalletAlreadyExist   = errors.New("wallet already exist")
	ErrWalletNotFound       = errors.New("wallet not found")
	ErrCardNotFound         = errors.New("card not found")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrCardLimitExceeded    = errors.New("card limit exceeded")
	ErrNotEnoughBalance     = errors.New("not enough balance")
)
