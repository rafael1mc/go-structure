package main

import (
	"gomodel/cmd/migration/migration"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
)

func Teardown() {

}

func Init(
	migration *migration.Migration,
	logger *slog.Logger,
) error {
	err := migration.Run()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	logger.Info("Migrations completed.")
	return nil
}
