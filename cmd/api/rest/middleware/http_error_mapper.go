package middleware

import (
	"database/sql"
	"gomodel/internal/auth"
	"gomodel/internal/member"
	"gomodel/internal/user"
	"net/http"
)

func NewErrorStatusCodeMapper() map[error]int {
	return map[error]int{
		sql.ErrNoRows: http.StatusNotFound,

		auth.UserNotEnabledError:  http.StatusUnauthorized,
		auth.InvalidPasswordError: http.StatusUnauthorized,

		user.InvalidEmailError:     http.StatusBadRequest,
		user.InvalidPasswordError:  http.StatusBadRequest,
		user.InvalidInviterIDError: http.StatusBadRequest,

		member.InvalidInstitutionError:  http.StatusBadRequest,
		member.InvalidEmailError:        http.StatusBadRequest,
		member.InvalidCategoryError:     http.StatusBadRequest,
		member.InvalidInviterIDError:    http.StatusBadRequest,
		member.InvalidInviteStatusError: http.StatusBadRequest,
		member.ExpiredInviteCodeError:   http.StatusUnauthorized,
	}
}
