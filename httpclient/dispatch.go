package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Dispatch an HTTP request and (optionally) unmarshal the payload.
func Dispatch(client *http.Client, r *http.Request, cha Challenger, v any) error {
	var (
		retry *http.Request
		err   error
	)

	if cha != nil {
		retry, err = cloneRequest(r)
		if err != nil {
			return err
		}
	}

	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf("Body close error: %v", err)
		}
	}()

	if cha != nil && cha.ChallengeAccepted(resp) {
		retry, err = cha.ChallengeResponse(retry, resp)
		if err != nil {
			return err
		}

		return Dispatch(client, retry, nil, v)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return errRequestUnauthorised
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if v != nil {
		return json.Unmarshal(b, v)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("%w: %s: status %d (body: %s)", errRequestUnsuccessful, r.URL.Path, resp.StatusCode, b)
	}

	return nil
}
