package migration

import (
	"gomodel/internal/shared/env"
	"gomodel/internal/shared/util/consts"
	"log/slog"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

type Migration struct {
	m      *migrate.Migrate
	logger *slog.Logger
}

func NewMigration(
	env *env.Env,
	db *sqlx.DB,
	logger *slog.Logger,
) *Migration {
	migration := &Migration{
		logger: logger,
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{DatabaseName: env.Database.Name})
	if err != nil {
		logger.Error("Failed to get database driver.", consts.SlogError(err))
		return migration
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		env.Database.Name,
		driver,
	)
	if err != nil {
		logger.Error("Failed to read migrations.", consts.SlogError(err))
		return migration
	}

	migration.m = m
	return migration
}

func (m Migration) Run() error {
	if m.m == nil {
		return MigrationNotInitializedError
	}

	return m.m.Up()
}
