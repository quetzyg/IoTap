package httpclient

import (
	"bytes"
	"io"
	"net/http"
)

// Challenger defines methods to handle HTTP authentication challenges.
// Implementations determine if a response requires authentication
// and generate the appropriate challenge response for subsequent requests.
type Challenger interface {
	ChallengeAccepted(*http.Response) bool
	ChallengeResponse(r *http.Request, resp *http.Response) (*http.Request, error)
}

// cloneRequest creates a replica of an *http.Request, including its body.
func cloneRequest(r *http.Request) (*http.Request, error) {
	clone := r.Clone(r.Context())

	if r.Body != nil {
		var body bytes.Buffer
		if _, err := io.Copy(&body, r.Body); err != nil {
			return nil, err
		}

		// Reset the original request body to allow re-use
		r.Body = io.NopCloser(bytes.NewBuffer(body.Bytes()))

		// Set the cloned request body
		clone.Body = io.NopCloser(bytes.NewBuffer(body.Bytes()))
	}

	return clone, nil
}
