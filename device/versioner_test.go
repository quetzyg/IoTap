package device

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type versioner struct {
	resource
	funcError error
}

func (v *versioner) VersionRequest() (*http.Request, error) {
	if v.funcError != nil {
		return nil, v.funcError
	}

	return http.NewRequest(http.MethodGet, "", nil)
}

func (v *versioner) Outdated() bool {
	return false
}
func (v *versioner) UpdateDetails() string {
	return ""
}

func TestVersion(t *testing.T) {
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
			dev: &versioner{
				funcError: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "failure: http response error",
			dev:  &versioner{},
			roundTripper: &RoundTripper{
				err: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "success",
			dev:  &versioner{},
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

			Version(tap, test.dev, ch)

			result := <-ch

			if _, ok := test.err.(*url.Error); ok {
				var urlErr *url.Error
				if !errors.As(result.err, &urlErr) {
					t.Fatalf("expected %#v, got %#v", test.err, result.err)
				}
				return
			}

			if !errors.Is(result.err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, result.err)
			}
		})
	}
}
