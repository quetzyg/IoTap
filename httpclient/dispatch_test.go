package httpclient

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// RoundTripper is a custom type that allows creating a mock RoundTripper
type RoundTripper struct {
	response *http.Response
	err      error
}

// RoundTrip implements the RoundTripper interface
func (rt RoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt.response, rt.err
}

// badReader returns an error when read.
type badReader struct{}

// Read implements the io.Reader interface.
func (badReader) Read(_ []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

// Close implements the io.Closer interface.
func (badReader) Close() error {
	return nil
}

var (
	uri = &url.URL{
		Scheme: "http",
		Host:   "192.168.146.12",
		Path:   "/settings",
	}

	v map[string]any
)

func TestDispatch(t *testing.T) {
	tests := []struct {
		name         string
		req          *http.Request
		roundTripper http.RoundTripper
		v            any
		err          error
	}{
		{
			name: "failure: timeout",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    uri,
			},
			roundTripper: &RoundTripper{
				response: nil,
				err:      context.DeadlineExceeded,
			},
			err: context.DeadlineExceeded,
		},

		{
			name: "failure: unable to read body",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    uri,
			},
			roundTripper: &RoundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &badReader{},
				},
			},
			err: io.ErrUnexpectedEOF,
		},

		{
			name: "failure: bad request",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    uri,
			},
			roundTripper: &RoundTripper{
				response: &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(strings.NewReader("Bad Request")),
				},
			},

			err: errResponse,
		},

		{
			name: "success: no body",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    uri,
			},
			roundTripper: &RoundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
				},
			},
		},

		{
			name: "success: unmarshal body",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    uri,
			},
			roundTripper: &RoundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
			v: &v,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &http.Client{
				Transport: test.roundTripper,
			}

			err := Dispatch(client, test.req, test.v)

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %v, got %v", test.err, err)
			}
		})
	}
}
