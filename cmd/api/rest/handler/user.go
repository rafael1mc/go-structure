package handler

import (
	"encoding/json"
	"gomodel/internal/user"

	"github.com/labstack/echo/v4"
)

func (h Handler) HandleCreateUser(c echo.Context) error {
	var createUserRequest user.CreateUserRequest
	err := json.NewDecoder(c.Request().Body).Decode(&createUserRequest)
	if err != nil {
		return err
	}

	_, err = h.user.Create(createUserRequest, nil)
	if err != nil {
		return err
	}

	return nil
}
