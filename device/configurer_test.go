package device

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type configurer struct {
	funcError error
	resource
}

func (c *configurer) ConfigureRequests(Config) ([]*http.Request, error) {
	if c.funcError != nil {
		return nil, c.funcError
	}

	return []*http.Request{
		{
			URL:    &url.URL{},
			Method: http.MethodGet,
		},
	}, nil
}

type configChallenger struct {
	configurer
}

func (cc *configChallenger) ChallengeAccepted(*http.Response) bool {
	return true
}

func (cc *configChallenger) ChallengeResponse(r *http.Request, _ *http.Response) (*http.Request, error) {
	return r, nil
}

func TestConfigure(t *testing.T) {
	tests := []struct {
		rt   http.RoundTripper
		dev  Resource
		err  error
		name string
	}{
		{
			name: "failure: unsupported procedure",
			dev:  &resource{},
			err:  ErrUnsupportedProcedure,
		},
		{
			name: "failure: function error",
			dev: &configurer{
				funcError: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "failure: http response error",
			dev:  &configurer{},
			rt: &roundTripper{
				err: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "success",
			dev:  &configurer{},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
		},
		{
			name: "success: challenger implementation",
			dev:  &configChallenger{},
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

			Configure(tap, test.dev, ch)

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
