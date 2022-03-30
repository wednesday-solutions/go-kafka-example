package producerclient

import (
	"fmt"
	"github.com/labstack/echo"
	client "libs/restclient"
	"net/http"
	"os"
)

type PongResponse struct {
	Ok string `json:"ok"`
}

// What handles GET requests to the /ping-what endpoint of the producer service
func What(ctx echo.Context) error {
	producerServiceURL := os.Getenv("PRODUCER_SVC_ENDPOINT") + "/ping-what"
	var pongJson PongResponse
	err := client.Get(producerServiceURL, &pongJson)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return ctx.JSON(http.StatusOK, pongJson)
}
