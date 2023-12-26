package seed

import (
	"fmt"
	"gomodel/internal/shared/database"
	"log/slog"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

type Seed struct {
	database *database.Database
	logger   *slog.Logger
}

func NewSeed(
	database *database.Database,
	logger *slog.Logger,
) *Seed {
	return &Seed{
		logger:   logger,
		database: database,
	}
}

func (s Seed) Run() error {
	seedDir, err := os.Open("seed/")
	if err != nil {
		return err
	}
	dirNames, err := seedDir.Readdirnames(0)
	if err != nil {
		return err
	}

	var tx *sqlx.Tx
	myTx, err := s.database.WithTx(&tx)
	if err != nil {
		return err
	}
	defer myTx.HandleRollback(&err)

	for _, dirName := range dirNames {
		fileContent, err := os.ReadFile(fmt.Sprintf("seed/%s", dirName))
		if err != nil {
			return err
		}
		_, err = tx.Query(string(fileContent))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
