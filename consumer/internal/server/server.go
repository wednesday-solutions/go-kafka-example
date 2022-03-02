package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
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
	e.GET("yo", what)
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

func what(ctx echo.Context) error {
	resp, err := http.Get("http://localhost:9001/say-what")
	if err != nil {

	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	data := string(body)
	return ctx.JSON(http.StatusOK, response{
		data,
	})
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
