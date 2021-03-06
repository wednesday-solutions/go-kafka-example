package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	gotemplate "consumer"
	"consumer/daos"
	"consumer/graphql_models"
	"consumer/internal/middleware/auth"
	"consumer/models"
	"consumer/pkg/utl/convert"
	resultwrapper "consumer/pkg/utl/result_wrapper"
)

func (r *mutationResolver) CreateRole(
	ctx context.Context,
	input graphql_models.RoleCreateInput) (*graphql_models.RolePayload, error) {

	userID := auth.UserIDFromContext(ctx)
	user, err := daos.FindUserByID(userID)
	if err != nil {
		return &graphql_models.RolePayload{}, resultwrapper.ResolverSQLError(err, "data")
	}
	userRole, err := daos.FindRoleByID(convert.NullDotIntToInt(user.RoleID))
	if err != nil {
		return &graphql_models.RolePayload{}, resultwrapper.ResolverSQLError(err, "data")
	}
	role := models.Role{
		AccessLevel: input.AccessLevel,
		Name:        input.Name,
	}
	if userRole.AccessLevel != int(gotemplate.SuperAdminRole) {
		return &graphql_models.RolePayload{}, fmt.Errorf("you don't appear to have enough access level for this request ")
	}

	newRole, err := daos.CreateRoleTx(role, nil)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "role")
	}
	return &graphql_models.RolePayload{Role: &graphql_models.Role{
		AccessLevel: newRole.AccessLevel,
		Name:        newRole.Name,
	},
	}, err
}
