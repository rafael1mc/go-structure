package handler

import (
	"encoding/json"
	"gomodel/cmd/api/rest/model"
	"gomodel/internal/member"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetMembers(c echo.Context) error {
	institutionID := c.Param("id")

	res, err := h.member.GetMembersByInstitutionID(institutionID, nil)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.GetListResponse(res))
}

func (h Handler) InviteMembers(c echo.Context) error {
	var inviteMemberRequest member.InviteMemberRequest
	err := json.NewDecoder(c.Request().Body).Decode(&inviteMemberRequest)
	if err != nil {
		return err
	}

	inviteMemberRequest.InstitutionID = c.Param("id")
	inviteMemberRequest.InviterID = c.Get("user_id").(string)

	err = h.member.InviteMember(inviteMemberRequest, nil)
	if err != nil {
		return err
	}

	return nil
}

func (h Handler) ValidateInvite(c echo.Context) error {
	var request member.ValidateInviteRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return err
	}

	request.InstitutionID = c.Param("id")

	_, err = h.member.ValidateInvite(request, nil)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h Handler) AcceptInvite(c echo.Context) error {
	var request member.AcceptInviteRequest
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return err
	}

	request.InstitutionID = c.Param("id")

	err = h.member.AcceptInvite(request, nil)
	if err != nil {
		return err
	}

	return nil
}
