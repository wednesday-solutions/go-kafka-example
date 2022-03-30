package restclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func Get(url string, jsonResponse interface{}) error {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("restclient.Get: failed to initialize request: %w", err)
	}
	response, err := Client.Do(request)
	if err != nil {
		return fmt.Errorf("restclient.Get: no response: %w", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("restclient.Get: failed to read response body: %w", err)
	}
	unmarshalErr := json.Unmarshal(body, &jsonResponse)
	if unmarshalErr != nil {
		return fmt.Errorf("restclient.Get: failed to read response body: %w", unmarshalErr)
	}
	return nil
}
