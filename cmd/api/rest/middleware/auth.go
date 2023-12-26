package middleware

import (
	"fmt"
	"gomodel/internal/auth"
	"gomodel/internal/shared/redis"
	"log/slog"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	redisWrapper *redis.RedisWrapper
	auth         *auth.Auth
	logger       *slog.Logger
}

func NewAuthMiddleware(
	redisWrapper *redis.RedisWrapper,
	auth *auth.Auth,
	logger *slog.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		redisWrapper: redisWrapper,
		auth:         auth,
		logger:       logger,
	}
}

func (m *AuthMiddleware) Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			route := c.Path()
			fmt.Println(route)
			if m.shouldSkipAuth(route) {
				return next(c)
			}

			authHeaders := c.Request().Header["Authorization"]
			if len(authHeaders) == 0 {
				return fmt.Errorf("unauthorized")
			}

			accessToken := strings.ReplaceAll(authHeaders[0], "Bearer ", "")
			if accessToken == "" {
				return fmt.Errorf("unauthorized")
			}

			userID, err := m.auth.ValidateAccessToken(accessToken)
			if err != nil {
				return err
			}

			c.Set("user_id", userID)

			return next(c)
		}
	}
}

func (m *AuthMiddleware) shouldSkipAuth(route string) bool {
	return route == "/ping" ||
		route == "/auth" ||
		route == "/auth/refresh" ||
		route == "/user" ||
		route == "/institution/:id/validate-invite" ||
		route == "/institution/:id/accept-invite"
}
