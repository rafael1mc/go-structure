package rest

import (
	"fmt"
	"gomodel/cmd/api/rest/handler"
	"gomodel/cmd/api/rest/middleware"
	"gomodel/internal/shared/env"
	"log/slog"

	"github.com/labstack/echo/v4"
	echomid "github.com/labstack/echo/v4/middleware"
)

type RestServer struct {
	handler          *handler.Handler
	logger           *slog.Logger
	env              *env.Env
	authMiddleware   *middleware.AuthMiddleware
	httpErrorHandler *middleware.HttpErrorHandler
}

func NewRestServer(
	handler *handler.Handler,
	logger *slog.Logger,
	env *env.Env,
	authMiddleware *middleware.AuthMiddleware,
	httpErrorHandler *middleware.HttpErrorHandler,
) *RestServer {
	return &RestServer{
		handler:          handler,
		logger:           logger,
		env:              env,
		authMiddleware:   authMiddleware,
		httpErrorHandler: httpErrorHandler,
	}
}

func (r RestServer) Run() {
	// Echo instance
	r.logger.Info("Initializing rest server...")
	e := echo.New()
	e.HideBanner = true
	// e.HidePort = true
	e.Debug = r.env.Environment.IsDebug
	e.HTTPErrorHandler = r.httpErrorHandler.Handler

	// Middleware
	r.logger.Debug("Setup middlewares")
	e.Use(echomid.Logger()) // TODO maybe the logger middleware should be using the project logger instead
	if !r.env.Environment.IsDebug {
		e.Use(echomid.Recover())
	}
	e.Use(r.authMiddleware.Auth())

	// Routes
	r.logger.Debug("Setup routes")
	e.GET("/ping", r.handler.HandlePing)
	e.POST("/auth", r.handler.HandleAuth)
	e.POST("/auth/refresh", r.handler.HandleRefreshAccessToken)
	e.POST("/user", r.handler.HandleCreateUser)
	e.PUT("/user/logout", r.handler.HandleLogout)

	e.GET("/institutions", r.handler.GetMyInstitution)
	e.GET("/institution/:id", r.handler.GetInstitution)

	e.GET("/institution/:id/members", r.handler.GetMembers)
	e.PUT("/institution/:id/invite-member", r.handler.InviteMembers)
	e.POST("/institution/:id/validate-invite", r.handler.ValidateInvite)
	e.POST("/institution/:id/accept-invite", r.handler.AcceptInvite)

	// Start server
	port := fmt.Sprintf(":%d", r.env.Api.Port)
	r.logger.Info("Starting HTTP Server", slog.Int("port", r.env.Api.Port))
	err := e.Start(port)
	e.Logger.Fatal(err)
}
