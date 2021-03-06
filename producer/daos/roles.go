package daos

import (
	"context"
	"database/sql"

	"producer/models"

	"github.com/volatiletech/sqlboiler/boil"
)

// CreateRoleTx ...
func CreateRoleTx(role models.Role, tx *sql.Tx) (models.Role, error) {
	contextExecutor := getContextExecutor(tx)

	err := role.Insert(context.Background(), contextExecutor, boil.Infer())
	return role, err
}

// FindRoleByID ...
func FindRoleByID(roleID int) (*models.Role, error) {
	contextExecutor := getContextExecutor(nil)
	return models.FindRole(context.Background(), contextExecutor, roleID)
}
