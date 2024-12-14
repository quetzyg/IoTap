package device

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type deployer struct {
	resource
	funcError error
}

func (d *deployer) DeployRequests(*http.Client, []*Script) ([]*http.Request, error) {
	if d.funcError != nil {
		return nil, d.funcError
	}

	return []*http.Request{
		{
			URL:    &url.URL{},
			Method: http.MethodGet,
		},
	}, nil
}

func TestDeploy(t *testing.T) {
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
			dev: &deployer{
				funcError: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "failure: http response error",
			dev:  &deployer{},
			roundTripper: &RoundTripper{
				err: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "success",
			dev:  &deployer{},
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

			Deploy(tap, test.dev, ch)

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
