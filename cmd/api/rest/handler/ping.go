package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HandlePing(c echo.Context) error {
	return c.JSON(http.StatusOK, h.ping.Respond())
}
