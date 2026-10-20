package errors

import "errors"

var (
	ErrCodeExpiredOrNotFound = errors.New("code not found or expired")
)
