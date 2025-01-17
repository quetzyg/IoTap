package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Dispatcher handles HTTP request dispatching with support for functional options.
type Dispatcher struct {
	client     *http.Client
	challenger Challenger
	bind       any
}

// DispatchOption is a function that modifies dispatch behavior.
type DispatchOption func(*Dispatcher)

// WithBinding returns a DispatchOption that configures response binding to the
// provided target. The target must be a pointer to a type that can receive the
// response data.
func WithBinding(bind any) DispatchOption {
	return func(d *Dispatcher) {
		d.bind = bind
	}
}

// WithChallenger returns a DispatchOption that adds challenge handling to the
// dispatch operation.
func WithChallenger(challenger Challenger) DispatchOption {
	return func(d *Dispatcher) {
		d.challenger = challenger
	}
}

// NewDispatcher creates a new *Dispatcher instance with the provided HTTP client.
func NewDispatcher(client *http.Client) *Dispatcher {
	return &Dispatcher{
		client: client,
	}
}

// Dispatch an HTTP request and (optionally) unmarshal the payload.
func (d *Dispatcher) Dispatch(r *http.Request, opts ...DispatchOption) error {
	for _, opt := range opts {
		opt(d)
	}

	var (
		retry *http.Request
		err   error
	)

	if d.challenger != nil {
		retry, err = cloneRequest(r)
		if err != nil {
			return err
		}
	}

	resp, err := d.client.Do(r)
	if err != nil {
		return err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf("Body close error: %v", err)
		}
	}()

	if d.challenger != nil && d.challenger.ChallengeAccepted(resp) {
		retry, err = d.challenger.ChallengeResponse(retry, resp)
		if err != nil {
			return err
		}

		return d.Dispatch(retry, WithChallenger(nil))
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return errRequestUnauthorised
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if d.bind != nil {
		return json.Unmarshal(b, &d.bind)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("%w: %s: status %d (body: %s)", errRequestUnsuccessful, r.URL.Path, resp.StatusCode, b)
	}

	return nil
}
