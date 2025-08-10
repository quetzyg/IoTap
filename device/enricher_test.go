package device

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type enricher struct {
	funcError error
	resource
}

func (e *enricher) EnrichRequest() (*http.Request, error) {
	if e.funcError != nil {
		return nil, e.funcError
	}

	return http.NewRequest(http.MethodGet, "", nil)
}

type enrichChallenger struct {
	enricher
}

func (ec *enrichChallenger) ChallengeAccepted(*http.Response) bool {
	return true
}

func (ec *enrichChallenger) ChallengeResponse(r *http.Request, _ *http.Response) (*http.Request, error) {
	return r, nil
}

func TestEnrich(t *testing.T) {
	tests := []struct {
		rt   http.RoundTripper
		dev  Resource
		err  error
		name string
	}{
		{
			name: "success: no enrichment required",
			dev:  &resource{},
			err:  nil,
		},
		{
			name: "failure: function error",
			dev: &enricher{
				funcError: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "failure: http response error",
			dev:  &enricher{},
			rt: &roundTripper{
				err: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "success",
			dev:  &enricher{},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
		},
		{
			name: "success: challenger implementation",
			dev:  &enrichChallenger{},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tap := &Tapper{
				transport: test.rt,
			}

			ch := make(chan *ProcedureResult, 1)

			Enrich(tap, test.dev, ch)

			result := <-ch

			var urlError *url.Error
			switch {
			case errors.As(test.err, &urlError):
				var ue *url.Error
				if errors.As(result.err, &ue) {
					return
				}

			case errors.Is(result.err, test.err):
				return

			default:
				t.Fatalf("expected %#v, got %#v", test.err, result.err)
			}
		})
	}
}
