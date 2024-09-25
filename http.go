package iotune

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Constants for setting the 'Content-Type' header to 'application/json' for all outgoing HTTP requests.
const (
	ContentTypeHeader = "Content-Type"
	JSONMimeType      = "application/json"
)

// Dispatch an HTTP request and (optionally) unmarshal the payload.
func Dispatch(client *http.Client, r *http.Request, v any) (err error) {
	response, err := client.Do(r)
	if err != nil {
		return
	}

	defer func() {
		err = errors.Join(err, response.Body.Close())
	}()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	if v != nil {
		return json.Unmarshal(b, v)
	}

	// From this point on, we only care about the HTTP request status
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("%s - %s (%d)", r.URL.Path, b, response.StatusCode)
	}

	return nil
}
