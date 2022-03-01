package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"consumer/daos"
	"consumer/graphql_models"
	"consumer/internal/middleware/auth"
	"consumer/pkg/utl/convert"
	resultwrapper "consumer/pkg/utl/result_wrapper"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (r *queryResolver) Me(ctx context.Context) (*graphql_models.User, error) {
	userID := auth.UserIDFromContext(ctx)
	user, err := daos.FindUserByID(userID)
	if err != nil {
		return &graphql_models.User{}, resultwrapper.ResolverSQLError(err, "data")
	}
	return convert.UserToGraphQlUser(user), err
}

func (r *queryResolver) Users(
	ctx context.Context,
	pagination *graphql_models.UserPagination) (*graphql_models.UsersPayload, error) {
	var queryMods []qm.QueryMod
	if pagination != nil {
		if pagination.Limit != 0 {
			queryMods = append(queryMods, qm.Limit(pagination.Limit), qm.Offset(pagination.Page*pagination.Limit))
		}
	}
	users, count, err := daos.FindAllUsersWithCount(queryMods)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	return &graphql_models.UsersPayload{Total: int(count), Users: convert.UsersToGraphQlUsers(users)}, nil
}

// Query returns graphql_models.QueryResolver implementation.
func (r *Resolver) Query() graphql_models.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
