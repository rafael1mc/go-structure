package auth

import "errors"

var (
	UserNotEnabledError  = errors.New("user is not enabled")
	InvalidPasswordError = errors.New("invalid password")
)
