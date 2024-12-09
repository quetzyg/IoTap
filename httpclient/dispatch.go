package httpclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

var errResponse = errors.New("HTTP response error")

// Dispatch an HTTP request and (optionally) unmarshal the payload.
func Dispatch(client *http.Client, r *http.Request, v any) (err error) {
	response, err := client.Do(r)
	if err != nil {
		return
	}

	defer func() {
		err = response.Body.Close()
		if err != nil {
			log.Printf("Body close error: %v", err)
		}
	}()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	if v != nil {
		return json.Unmarshal(b, v)
	}

	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("%w: %s: status %d (body: %s)", errResponse, r.URL.Path, response.StatusCode, b)
	}

	return nil
}
