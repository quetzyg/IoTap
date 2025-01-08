package device

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type securer struct {
	resource
	funcError error
}

func (s *securer) Secured() bool { return true }

func (s *securer) SetCredentials(_ *Credentials) {}

func (s *securer) AuthConfigRequest(*AuthConfig) (*http.Request, error) {
	if s.funcError != nil {
		return nil, s.funcError
	}

	return http.NewRequest(http.MethodGet, "", nil)

}

func TestSecure(t *testing.T) {
	tests := []struct {
		name         string
		roundTripper http.RoundTripper
		dev          Resource
		err          error
	}{
		{
			name: "failure: unsupported procedure",
			dev:  &resource{},
			err:  ErrUnsupportedProcedure,
		},
		{
			name: "failure: function error",
			dev: &securer{
				funcError: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "failure: http response error",
			dev:  &securer{},
			roundTripper: &RoundTripper{
				err: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "success",
			dev:  &securer{},
			roundTripper: &RoundTripper{
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
				transport: test.roundTripper,
			}

			ch := make(chan *ProcedureResult, 1)

			Secure(tap, test.dev, ch)

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
