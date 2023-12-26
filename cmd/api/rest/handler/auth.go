package handler

import (
	"encoding/json"
	"gomodel/internal/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) HandleAuth(c echo.Context) error {
	var authRequest auth.AuthRequest
	err := json.NewDecoder(c.Request().Body).Decode(&authRequest)
	if err != nil {
		return err
	}

	jwt, err := h.auth.AuthenticateUser(authRequest)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, jwt)
}

func (h Handler) HandleRefreshAccessToken(c echo.Context) error {
	var refreshTokenRequest auth.RefreshTokenRequest
	err := json.NewDecoder(c.Request().Body).Decode(&refreshTokenRequest)
	if err != nil {
		return err
	}

	jwt, err := h.auth.RefreshAccessToken(refreshTokenRequest)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, jwt)
}

func (h Handler) HandleLogout(c echo.Context) error {
	userID := c.Get("user_id").(string)

	return h.auth.Logout(userID)
}
