package producerclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	client "libs/restclient"
	"net/http"
	"os"

	"github.com/machinebox/graphql"
)

type PongResponse struct {
	Ok string `json:"ok"`
}

// What handles GET requests to the /ping endpoint of the consumer service
func What(ctx echo.Context) error {
	producerServiceURL := os.Getenv("PRODUCER_SVC_ENDPOINT") + "/ping-what"
	var pongJson PongResponse
	err := client.Get(producerServiceURL, &pongJson)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return ctx.JSON(http.StatusOK, pongJson)
}

type LoginResponse struct {
	Token string `json:"token"`
}

type TokenResponse struct {
	Login LoginResponse `json:"login"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginPassthrough handles POST requests to the /producer-login endpoint of the consumer service
func LoginPassthrough(ctx echo.Context) error {
	producerServiceURL := os.Getenv("PRODUCER_SVC_ENDPOINT") + "/graphql"
	graphqlClient := graphql.NewClient(producerServiceURL)

	req := graphql.NewRequest(`
    mutation login($username: String!, $password: String!) {
        login (username: $username, password: $password) {
            token
        }
    }
`)
	var credentials Credentials
	if err := json.NewDecoder(ctx.Request().Body).Decode(&credentials); err != nil {
		return fmt.Errorf("LoginPassthrough::failed to decode request body::%w", err)
	}

	req.Var("username", credentials.Username)
	req.Var("password", credentials.Password)

	var tokenResponse TokenResponse
	if err := graphqlClient.Run(context.Background(), req, &tokenResponse); err != nil {
		return fmt.Errorf("%w", err)
	}
	return ctx.JSON(http.StatusOK, &tokenResponse)
}
