package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"producer/daos"
	"producer/graphql_models"
	"producer/internal/config"
	"producer/internal/middleware/auth"
	"producer/internal/service"
	"producer/models"
	"producer/pkg/utl/convert"
	kafka "producer/pkg/utl/kafkaservice"
	resultwrapper "producer/pkg/utl/result_wrapper"

	"github.com/volatiletech/null"
)

func (r *mutationResolver) CreateUser(
	ctx context.Context,
	input graphql_models.UserCreateInput) (*graphql_models.UserPayload, error) {

	user := models.User{
		Username:  null.StringFromPtr(input.Username),
		Password:  null.StringFromPtr(input.Password),
		Email:     null.StringFromPtr(input.Email),
		FirstName: null.StringFromPtr(input.FirstName),
		LastName:  null.StringFromPtr(input.LastName),
		RoleID:    convert.PointerStringToNullDotInt(input.RoleID),
	}
	// loading configurations
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("error in loading config ")
	}
	// creating new secure service
	sec := service.Secure(cfg)
	user.Password = null.StringFrom(sec.Hash(user.Password.String))
	newUser, err := daos.CreateUserTx(user, nil)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "user information")
	}
	graphUser := convert.UserToGraphQlUser(&newUser)

	r.Lock()
	for _, observer := range r.Observers {
		observer <- graphUser
	}

	b, _ := json.Marshal(newUser)
	r.Unlock()
	go kafka.Produce(
		context.Background(),
		kafka.NEW_USER_CREATED,
		[]byte(strconv.Itoa(rand.Int())),
		b,
	)
	return &graphql_models.UserPayload{User: graphUser}, err
}

func (r *mutationResolver) UpdateUser(
	ctx context.Context,
	input *graphql_models.UserUpdateInput) (*graphql_models.UserUpdatePayload, error) {

	userID := auth.UserIDFromContext(ctx)
	u := models.User{
		ID:        userID,
		FirstName: null.StringFromPtr(input.FirstName),
		LastName:  null.StringFromPtr(input.LastName),
		Mobile:    null.StringFromPtr(input.Mobile),
		Phone:     null.StringFromPtr(input.Phone),
		Address:   null.StringFromPtr(input.Address),
	}
	_, err := daos.UpdateUserTx(u, nil)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "new information")
	}

	graphUser := convert.UserToGraphQlUser(&u)
	r.Lock()
	for _, observer := range r.Observers {
		observer <- graphUser
	}
	r.Unlock()

	return &graphql_models.UserUpdatePayload{Ok: true}, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context) (*graphql_models.UserDeletePayload, error) {
	userID := auth.UserIDFromContext(ctx)
	u, err := daos.FindUserByID(userID)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "data")
	}
	_, err = daos.DeleteUser(*u)
	if err != nil {
		return nil, resultwrapper.ResolverSQLError(err, "user")
	}
	return &graphql_models.UserDeletePayload{ID: fmt.Sprint(userID)}, nil
}
