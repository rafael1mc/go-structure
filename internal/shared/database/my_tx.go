package database

import (
	"database/sql"
	"fmt"
	"gomodel/internal/shared/env"
	"gomodel/internal/shared/util/consts"
	"log/slog"
	"runtime"

	"github.com/jmoiron/sqlx"
)

type MyTx struct {
	database     *Database
	tx           *sqlx.Tx
	isExternalTx bool
	logger       *slog.Logger
	env          *env.Env
}

func (t MyTx) IsExternal() bool {
	return t.isExternalTx
}

func (t MyTx) IsInternal() bool {
	return !t.IsExternal()
}

func (t MyTx) HandleRollback(inputErr *error) {
	if t.isExternalTx {
		fmt.Println("External TX, skipping rollback.")
		return
	}
	if *inputErr != nil {
		fmt.Println("Internal TX with error, running rollback.")
		t.tx.Rollback()
		return
	}

	_, err := t.tx.Exec(`SELECT 1`)
	if err == sql.ErrTxDone {
		// transaction done successfully
		return
	}

	// if reached this point, commit/rollback was not called inside func that created tx

	var callerName string
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		callerName = details.Name()
	}

	err = fmt.Errorf("leaking transaction. Forgot to call tx.Commit? Maybe at %s", callerName)

	if t.env.Environment.Name != consts.ProductionEnv {
		// crash app to be more direct when not in prod
		panic(err)
	} else {
		t.logger.Error("error at HandleRollback", consts.SlogError(err))
	}
}
