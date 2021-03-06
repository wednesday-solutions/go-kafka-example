package daos

import (
	"database/sql"

	"github.com/volatiletech/sqlboiler/boil"
)

func getContextExecutor(tx *sql.Tx) (contextExecutor boil.ContextExecutor) {
	if tx == nil {
		contextExecutor = boil.GetContextDB()
	} else {
		contextExecutor = boil.ContextExecutor(tx)
	}
	return contextExecutor
}
