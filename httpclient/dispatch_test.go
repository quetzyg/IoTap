package httpclient

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// roundTripper is a custom type used for mocking HTTP responses.
type roundTripper struct {
	response *http.Response
	err      error
}

// RoundTrip implements the http.RoundTripper interface.
func (rt *roundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
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

type challenger struct {
	req *http.Request
	err error
}

func (c *challenger) ChallengeAccepted(resp *http.Response) bool {
	return resp.StatusCode == http.StatusUnauthorized
}

func (c *challenger) ChallengeResponse(_ *http.Request, _ *http.Response) (*http.Request, error) {
	if c.err != nil {
		return nil, c.err
	}

	return c.req, nil
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
		name       string
		req        *http.Request
		rt         http.RoundTripper
		challenger Challenger
		v          any
		err        error
	}{
		{
			name: "failure: clone request error",
			req: &http.Request{
				Body: &badReader{},
			},
			challenger: &challenger{},
			err:        io.ErrUnexpectedEOF,
		},
		{
			name: "failure: timeout",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    uri,
			},
			rt: &roundTripper{
				response: nil,
				err:      net.ErrClosed,
			},
			err: net.ErrClosed,
		},
		{
			name: "failure: challenge accepted error",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    uri,
			},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusUnauthorized,
				},
			},
			challenger: &challenger{
				err: net.ErrClosed,
			},
			err: net.ErrClosed,
		},
		{
			name: "failure: challenge unauthorised",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    uri,
			},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusUnauthorized,
				},
			},
			challenger: &challenger{
				req: &http.Request{
					Method: http.MethodGet,
					URL:    uri,
					Header: http.Header{
						AuthorizationHeader: []string{`Digest foobar`},
					},
				},
			},
			err: errRequestUnauthorised,
		},
		{
			name: "failure: unable to read body",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    uri,
			},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &badReader{},
				},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			name: "success: unmarshal body",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    uri,
			},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
			v: &v,
		}, {
			name: "failure: bad request",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    uri,
			},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusBadRequest,
					Body:       io.NopCloser(strings.NewReader("Bad Request")),
				},
			},

			err: errRequestUnsuccessful,
		},
		{
			name: "success: no body",
			req: &http.Request{
				Method: http.MethodGet,
				URL:    uri,
			},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &http.Client{
				Transport: test.rt,
			}

			err := Dispatch(client, test.req, test.challenger, test.v)

			var urlError *url.Error
			switch {
			case errors.As(test.err, &urlError):
				var ue *url.Error
				if errors.As(err, &ue) {
					return
				}

			default:
				if errors.Is(err, test.err) {
					return
				}
			}

			t.Fatalf("expected %#v, got %#v", test.err, err)
		})
	}
}
