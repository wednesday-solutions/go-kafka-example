package server_test

import (
	"bytes"
	"consumer/internal/server"
	"consumer/testutls"
	"context"
	"fmt"
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"libs/restclient"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"testing"
	"time"
	"utils/mocks"
)

func init() {
	restclient.Client = &mocks.MockClient{}
}

// Improve tests
func TestHealthCheck(t *testing.T) {
	e := server.New()
	if e == nil {
		t.Errorf("Server should not be nil")
	}
	assert.NotEmpty(t, e)
	response, err := testutls.MakeRequest(
		testutls.RequestParameters{
			E:          e,
			Pathname:   "/",
			HttpMethod: "GET",
		},
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, response["data"], "consumer: Go template at your service!🍲")
}

func TestProducerServiceApi(t *testing.T) {
	json := `{"ok": "pong"}`
	// create a new reader with that JSON
	responseBody := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       responseBody,
		}, nil
	}
	e := server.New()
	if e == nil {
		t.Errorf("Server should not be nil")
	}
	assert.NotEmpty(t, e)
	response, err := testutls.MakeRequest(
		testutls.RequestParameters{
			E:          e,
			Pathname:   "/ping",
			HttpMethod: "GET",
		},
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, "pong", response["ok"])
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("connection  error")
	}
	response1, _ := testutls.MakeRequest(
		testutls.RequestParameters{
			E:          e,
			Pathname:   "/ping",
			HttpMethod: "GET",
		},
	)
	fmt.Println(response1)
	assert.Equal(t, "Internal Server Error", response1["message"])
	assert.Empty(t, response1["ok"])
}

type args struct {
	e                    *echo.Echo
	cfg                  *server.Config
	startServer          func(e *echo.Echo, s *http.Server) (err error)
	startServerCalled    bool
	serverShutDownCalled bool
	shutDownFailed       bool
}

func initValues(shutDownFailed bool, startServer func(e *echo.Echo, s *http.Server) error) args {
	config := testutls.MockConfig()
	return args{
		e: server.New(),
		cfg: &server.Config{
			Port:                config.Server.Port,
			ReadTimeoutSeconds:  config.Server.ReadTimeout,
			WriteTimeoutSeconds: config.Server.WriteTimeout,
			Debug:               config.Server.Debug,
		},
		startServer:    startServer,
		shutDownFailed: true,
	}
}
func TestStart(t *testing.T) {

	cases := map[string]struct {
		args args
	}{
		"Success": {
			args: initValues(false, func(e *echo.Echo, s *http.Server) (err error) {
				return nil
			}),
		},
		"Failure_ServerStartUpFailed": {
			args: initValues(false, func(e *echo.Echo, s *http.Server) (err error) {
				return fmt.Errorf("error starting up")
			}),
		},
		"Failure_ServerShutDownFailed": {
			args: initValues(true, func(e *echo.Echo, s *http.Server) (err error) {
				return nil
			}),
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			ApplyMethod(reflect.TypeOf(tt.args.e), "StartServer", func(e *echo.Echo, s *http.Server) (err error) {
				err = tt.args.startServer(e, s)
				tt.args.startServerCalled = true
				return err
			})

			if tt.args.shutDownFailed {
				ApplyMethod(reflect.TypeOf(tt.args.e), "Shutdown", func(e *echo.Echo, ctx context.Context) (err error) {
					return fmt.Errorf("error shutting down")
				})

				ApplyMethod(
					reflect.TypeOf(tt.args.e.StdLogger), "Fatal",
					func(l *log.Logger, i ...interface{}) {
						tt.args.serverShutDownCalled = true
					})
			}

			go func() {
				time.Sleep(200 * time.Millisecond)
				proc, err := os.FindProcess(os.Getpid())
				if err != nil {
					log.Fatal(err)
				}
				sigc := make(chan os.Signal, 1)
				signal.Notify(sigc, os.Interrupt)
				go func() {
					<-sigc
					signal.Stop(sigc)
				}()
				err = proc.Signal(os.Interrupt)
				if err != nil {
					log.Fatal("errror")
				}
				time.Sleep(1 * time.Second)

			}()
			server.Start(tt.args.e, tt.args.cfg)
			time.Sleep(400 * time.Millisecond)
			assert.Equal(t, tt.args.startServerCalled, true)

			if tt.args.shutDownFailed {
				assert.Equal(t, tt.args.serverShutDownCalled, true)
			}

		})
	}
}
