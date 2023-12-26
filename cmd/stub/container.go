package main

import (
	"gomodel/internal/shared"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := shared.BuildSharedContainer()

	return container
}
