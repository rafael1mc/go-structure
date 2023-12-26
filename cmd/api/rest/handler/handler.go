package handler

import (
	"gomodel/internal/auth"
	"gomodel/internal/institution"
	"gomodel/internal/member"
	"gomodel/internal/ping"
	"gomodel/internal/user"
)

type Handler struct {
	ping        *ping.Ping
	auth        *auth.Auth
	user        *user.User
	institution *institution.Institution
	member      *member.Member
}

func NewHandler(
	ping *ping.Ping,
	auth *auth.Auth,
	user *user.User,
	institution *institution.Institution,
	member *member.Member,
) *Handler {
	return &Handler{
		ping:        ping,
		auth:        auth,
		user:        user,
		institution: institution,
		member:      member,
	}
}
