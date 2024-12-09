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

// RoundTripper is a custom type used for mocking HTTP responses.
type RoundTripper struct {
	response *http.Response
	err      error
}

// RoundTrip implements the RoundTripper interface.
func (rt RoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return rt.response, rt.err
}

func TestProbeIP(t *testing.T) {
	tests := []struct {
		name         string
		prober       Prober
		roundTripper http.RoundTripper
		res          Resource
		err          error
	}{
		{
			name:   "failure: bad prober",
			prober: &prober{badRequest: true},
			err:    errRequestCreate,
		},
		{
			name:   "failure: http response error",
			prober: &prober{},
			roundTripper: &RoundTripper{
				err: &url.Error{
					Op:  "parse",
					URL: ":",
					Err: errors.New("missing protocol scheme"),
				},
			},
		},
		{
			name:   "failure: unexpected device",
			prober: &prober{resource: &resource{unexpected: true}},
			roundTripper: &RoundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
		},
		{
			name:   "failure: parsing error",
			prober: &prober{resource: &resource{}},
			roundTripper: &RoundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("}")),
				},
			},
		},
		{
			name:   "success",
			prober: &prober{resource: &resource{}},
			roundTripper: &RoundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
			res: &resource{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &http.Client{
				Transport: test.roundTripper,
			}

			res, err := probeIP(test.prober, client, net.ParseIP("192.168.146.123"))

			if !reflect.DeepEqual(res, test.res) {
				t.Fatalf("expected %v, got %v", test.res, res)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %v, got %v", test.err, err)
			}
		})
	}
}
