package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"consumer/internal/middleware/secure"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// New instantiates new Echo server
func New() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover(), secure.CORS(), secure.Headers())
	e.GET("/", healthCheck)
	e.GET("ping", what)
	e.Validator = &CustomValidator{V: validator.New()}
	custErr := &customErrHandler{e: e}
	e.HTTPErrorHandler = custErr.handler
	e.Binder = &CustomBinder{b: &echo.DefaultBinder{}}
	return e
}

type response struct {
	Data string `json:"data"`
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, response{Data: "consumer: Go template at your service!🍲"})
}

type PongResponse struct {
	Ok string `json:"ok"`
}

func what(ctx echo.Context) error {
	producerServiceURL := url.URL{Scheme: "http", Host: os.Getenv("PRODUCER_SVC_ENDPOINT"), Path: "/ping-what"}
	fmt.Println(producerServiceURL.String())
	resp, err := http.Get(producerServiceURL.String())
	if err != nil {
		fmt.Println(err, "couldn't get a response")
	}
	defer resp.Body.Close()
	var m PongResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err, "couldn't read the body")
	}
	unmarshalErr := json.Unmarshal(body, &m)
	if unmarshalErr != nil {
		fmt.Println(unmarshalErr, "couldn't unmarshall the body")
	}
	fmt.Println(m, "unmarshalledResponse", body)
	return ctx.JSON(http.StatusOK, m)
}

// Config represents server specific config
type Config struct {
	Port                string
	ReadTimeoutSeconds  int
	WriteTimeoutSeconds int
	Debug               bool
}

// Start starts echo server
func Start(e *echo.Echo, cfg *Config) {
	s := &http.Server{
		Addr:         cfg.Port,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
	}
	e.Debug = cfg.Debug

	// Start server
	go func() {
		fmt.Print("Warming up server... ")
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.StdLogger.Fatal(err)
	}
}
