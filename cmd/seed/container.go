package main

import (
	"gomodel/cmd/seed/seed"
	"gomodel/internal/shared"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := shared.BuildSharedContainer()

	// CMD
	//   Seed
	container.Provide(seed.NewSeed)

	return container
}
