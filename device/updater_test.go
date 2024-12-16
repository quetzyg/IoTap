package device

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type updater struct {
	resource
	funcError error
}

func (u *updater) UpdateRequest() (*http.Request, error) {
	if u.funcError != nil {
		return nil, u.funcError
	}

	return http.NewRequest(http.MethodGet, "", nil)
}

func TestUpdate(t *testing.T) {
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
			dev: &updater{
				funcError: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "failure: http response error",
			dev:  &updater{},
			roundTripper: &RoundTripper{
				err: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "success",
			dev:  &updater{},
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

			Update(tap, test.dev, ch)

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
