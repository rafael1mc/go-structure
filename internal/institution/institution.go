package institution

import (
	"gomodel/internal/shared/database"
	"gomodel/internal/shared/timeprovider"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Institution struct {
	database     *database.Database
	logger       *slog.Logger
	timeProvider timeprovider.TimeProvider
}

func NewInstitution(
	database *database.Database,
	logger *slog.Logger,
	timeProvider timeprovider.TimeProvider,
) *Institution {
	return &Institution{
		database:     database,
		logger:       logger,
		timeProvider: timeProvider,
	}
}

func (i Institution) GetMyInstitutions(userID string, tx *sqlx.Tx) (res []InstitutionResponse, err error) {
	myTx, err := i.database.WithTx(&tx)
	if err != nil {
		return
	}
	defer myTx.HandleRollback(&err)

	err = tx.Select(
		&res,
		`SELECT i.* FROM institution i
		INNER JOIN member m ON c.institution_id = i.id
		INNER JOIN user_ u ON u.id = m.user_id
		WHERE m.is_enabled = TRUE
		AND u.is_enabled = TRUE
		AND u.id = $1`,
		userID,
	)

	return
}

func (i Institution) GetInstitutionByID(
	id string,
	tx *sqlx.Tx,
) (res InstitutionResponse, err error) {
	myTx, err := i.database.WithTx(&tx)
	if err != nil {
		return
	}
	defer myTx.HandleRollback(&err)

	err = tx.Get(
		&res,
		`SELECT * FROM institution i
		WHERE i.id = $1`,
		id,
	)

	return
}
