package handler

import (
	"gomodel/cmd/api/rest/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetMyInstitution(c echo.Context) error {
	userID := c.Get("user_id").(string)

	res, err := h.institution.GetMyInstitutions(userID, nil)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.GetListResponse(res))
}

func (h Handler) GetInstitution(c echo.Context) error {
	institutionID := c.Param("id")

	res, err := h.institution.GetInstitutionByID(institutionID, nil)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}
