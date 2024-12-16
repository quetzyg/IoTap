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

func TestNewTapper(t *testing.T) {
	tap := NewTapper([]Prober{
		&prober{},
	})

	if len(tap.probers) != 1 {
		t.Fatal("prober count must be 1")
	}
}

func TestTapper_SetConfig(t *testing.T) {
	tap := &Tapper{}

	tap.SetConfig(&config{Foo: "bar"})

	if tap.config.Empty() {
		t.Fatal("config is empty")
	}
}

func TestTapper_SetScripts(t *testing.T) {
	tap := &Tapper{}

	tap.SetScripts([]*Script{
		{
			name: "foo",
			code: []byte("var foo = 123;"),
		},
	})

	if len(tap.scripts) != 1 {
		t.Fatal("script count must be 1")
	}
}

func TestTapper_probe(t *testing.T) {
	tests := []struct {
		name         string
		prober       Prober
		roundTripper http.RoundTripper
		dev          Resource
		failed       bool
		err          error
	}{
		{
			name: "failure: probe error",
			prober: &prober{
				funcError: &url.Error{},
			},
			failed: true,
			err:    &ProbeError{},
		},
		{
			name:   "success: device found",
			prober: &prober{resource: &resource{}},
			roundTripper: &RoundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
			dev:    &resource{},
			failed: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tap := &Tapper{
				probers: []Prober{test.prober},
			}

			ch := make(chan *ProcedureResult, 1)

			client := &http.Client{
				Transport: test.roundTripper,
			}

			tap.probe(ch, client, net.ParseIP("192.168.146.123"))

			result := <-ch

			if test.failed != result.Failed() {
				t.Fatalf("expected %t, got %t", test.failed, result.Failed())
			}

			if !reflect.DeepEqual(result.dev, test.dev) {
				t.Fatalf("expected %#v, got %#v", test.dev, result.dev)
			}

			var probeError *ProbeError
			switch {
			case errors.As(test.err, &probeError):
				var pe *ProbeError
				if errors.As(result.err, &pe) {
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

// prober implementation for testing purposes.
type prober struct {
	resource  *resource
	funcError error
}

// Request implementation for testing purposes.
func (p *prober) Request(_ net.IP) (*http.Request, Resource, error) {
	if p.funcError != nil {
		return nil, nil, p.funcError
	}

	r, err := http.NewRequest(http.MethodGet, "", nil)

	return r, p.resource, err
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
			name: "failure: bad prober",
			prober: &prober{
				funcError: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name:   "failure: http response error",
			prober: &prober{},
			roundTripper: &RoundTripper{
				err: &url.Error{},
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
				t.Fatalf("expected %#v, got %#v", test.res, res)
			}

			var urlError *url.Error
			switch {
			case errors.As(test.err, &urlError):
				var ue *url.Error
				if errors.As(err, &ue) {
					return
				}

			case errors.Is(err, test.err):
				return

			default:
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}
