package device

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var errRequestCreate = errors.New("error creating request")

// prober implementation for testing purposes.
type prober struct {
	badRequest bool
	resource   *resource
}

// Request implementation for testing purposes.
func (p *prober) Request(_ net.IP) (*http.Request, Resource, error) {
	if p.badRequest {
		return nil, nil, errRequestCreate
	}

	r, _ := http.NewRequest(http.MethodGet, "http://192.168.146.123/settings", nil)

	return r, p.resource, nil
}

// RoundTripFunc is a custom type that allows creating a mock RoundTripper.
type RoundTripFunc func(req *http.Request) (*http.Response, error)

// RoundTrip implements the RoundTripper interface
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// newMockClient creates an HTTP client with a custom transport that can return errors
func newMockClient(mockFunc RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: mockFunc,
	}
}

func TestProbeIP(t *testing.T) {
	tests := []struct {
		name   string
		prober Prober
		client *http.Client
		res    Resource
		err    error
	}{
		{
			name:   "failure: bad prober",
			prober: &prober{badRequest: true},
			err:    errRequestCreate,
		},
		{
			name:   "failure: http response error",
			prober: &prober{},
			client: newMockClient(func(req *http.Request) (*http.Response, error) {
				return nil, &url.Error{
					Op:  "parse",
					URL: ":",
					Err: errors.New("missing protocol scheme"),
				}
			}),
		},
		{
			name:   "failure: unexpected device",
			prober: &prober{resource: &resource{unexpected: true}},
			client: newMockClient(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				}, nil
			}),
		},
		{
			name:   "failure: parsing error",
			prober: &prober{resource: &resource{}},
			client: newMockClient(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("}")),
				}, nil
			}),
		},
		{
			name:   "success",
			prober: &prober{resource: &resource{}},
			client: newMockClient(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				}, nil
			}),
			res: &resource{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := probeIP(test.prober, test.client, net.ParseIP("192.168.146.123"))

			if !reflect.DeepEqual(res, test.res) {
				t.Fatalf("expected %v, got %v", test.res, res)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %v, got %v", test.err, err)
			}
		})
	}
}
