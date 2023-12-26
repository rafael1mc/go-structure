package database

import (
	"fmt"
	"gomodel/internal/shared/env"
	"gomodel/internal/shared/util/consts"
	"log/slog"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	db     *sqlx.DB
	logger *slog.Logger
	env    *env.Env
}

func NewDatabase(
	env *env.Env,
	logger *slog.Logger,
) *Database {
	var db *sqlx.DB

	sslMode := ""
	if env.Environment.Name == consts.LocalEnv {
		sslMode = "sslmode=disable"
	}
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?%s",
		env.Database.User,
		env.Database.Pass,
		env.Database.Host,
		env.Database.Port,
		env.Database.Name,
		sslMode,
	)

	logger.Debug("Creating new Database Connection", slog.String("conn_string", databaseUrl))

	conn := func() error {
		var err error
		db, err = sqlx.Connect("postgres", databaseUrl)
		if err != nil {
			logger.Error("could not connect to database", consts.SlogError(err))
			return err
		}

		err = db.Ping()
		if err != nil {
			logger.Error("could not ping database", consts.SlogError(err))
			return err
		}

		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = 2 * time.Second
	expBackoff.MaxElapsedTime = 30 * time.Second
	expBackoff.Reset()

	err := backoff.Retry(conn, expBackoff)
	if err != nil {
		panic(fmt.Errorf("Failed to connect to database after retrying: %w", err))
	}

	return &Database{
		db:     db.Unsafe(), // Unsafe is used to no complain about missing destination names from query to mapped model
		logger: logger,
		env:    env,
	}
}

func (db *Database) WithTx(inputTx **sqlx.Tx) (myTx MyTx, err error) {
	isExternalTx := *inputTx != nil
	if !isExternalTx {
		var tx *sqlx.Tx
		tx, err = db.db.Beginx()
		if err != nil {
			return
		}
		*inputTx = tx
	}

	myTx = MyTx{
		database:     db,
		tx:           *inputTx,
		isExternalTx: isExternalTx,
		logger:       db.logger,
		env:          db.env,
	}
	return
}
