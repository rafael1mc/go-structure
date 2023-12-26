package main

import (
	"gomodel/cmd/migration/migration"
	"gomodel/internal/shared"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := shared.BuildSharedContainer()

	// CMD
	//   Migration
	container.Provide(migration.NewMigration)

	return container
}
