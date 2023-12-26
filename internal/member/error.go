package member

import "errors"

var (
	InvalidInstitutionError = errors.New("invalid institution")
	InvalidUserError        = errors.New("invalid user")
	InvalidEmailError       = errors.New("invalid email")
	InvalidCategoryError    = errors.New("invalid category")
	InvalidInviterIDError   = errors.New("invalid inviter")

	InvalidInviteStatusError = errors.New("invalid invite status")
	ExpiredInviteCodeError   = errors.New("expired invite code")

	InvalidNameError     = errors.New("invalid name")
	InvalidPasswordError = errors.New("invalid password")
)
