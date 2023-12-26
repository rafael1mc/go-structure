package main

import (
	"gomodel/cmd/api/rest"
	"gomodel/cmd/api/rest/handler"
	"gomodel/cmd/api/rest/middleware"
	"gomodel/internal/shared"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := shared.BuildSharedContainer()

	// CMD
	//   Rest
	container.Provide(middleware.NewAuthMiddleware)
	container.Provide(middleware.NewErrorStatusCodeMapper)
	container.Provide(middleware.NewHttpErrorHandler)
	container.Provide(rest.NewRestServer)
	container.Provide(handler.NewHandler)

	return container
}
