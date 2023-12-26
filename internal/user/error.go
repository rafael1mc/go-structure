package user

import "errors"

var (
	InvalidNameError      = errors.New("invalid name")
	InvalidEmailError     = errors.New("invalid email")
	InvalidPasswordError  = errors.New("invalid password")
	InvalidInviterIDError = errors.New("invalid inviter")
)
