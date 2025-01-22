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
	"time"
)

func TestNewTapper(t *testing.T) {
	tap := NewTapper(time.Second, []Prober{
		&prober{},
	})

	if len(tap.probers) != 1 {
		t.Fatal("prober count must be 1")
	}

	if tap.timeout != time.Second {
		t.Fatal("timeout must be 1 second")
	}
}

func TestTapper_SetCredentials(t *testing.T) {
	tap := &Tapper{}

	tap.SetCredentials(&Credentials{})

	if tap.cred == nil {
		t.Fatal("credentials are nil")

	}
}

func TestTapper_SetConfig(t *testing.T) {
	tap := &Tapper{}

	tap.SetConfig(&config{Foo: "bar"})

	if tap.config.Empty() {
		t.Fatal("config is empty")
	}
}

func TestTapper_SetAuthConfig(t *testing.T) {
	tap := &Tapper{}

	tap.SetAuthConfig(&AuthConfig{})

	if tap.auth == nil {
		t.Fatal("auth configuration is nil")
	}
}

func TestTapper_SetDeployment(t *testing.T) {
	tap := &Tapper{}

	tap.SetDeployment(&Deployment{
		Policy: &Policy{
			Mode: PolicyModeWhitelist,
		},
		Scripts: []*Script{
			{
				path: "/foo/bar.js",
				code: []byte("var foo = 123;"),
			},
		},
	})

	if len(tap.deployment.Scripts) != 1 {
		t.Fatal("script count must be 1")
	}
}

func TestTapper_probe(t *testing.T) {
	tests := []struct {
		name   string
		prober Prober
		rt     http.RoundTripper
		dev    Resource
		failed bool
		err    error
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
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
			dev:    &resource{},
			failed: false,
		},
		{
			name:   "success: secure device found",
			prober: &prober{resource: &securer{}},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
			dev:    &securer{},
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
				Transport: test.rt,
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
	resource  Resource
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

// roundTripper is a custom type used for mocking HTTP responses.
type roundTripper struct {
	response *http.Response
	err      error
}

// RoundTrip implements the http.RoundTripper interface.
func (rt *roundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return rt.response, rt.err
}

func TestProbeIP(t *testing.T) {
	tests := []struct {
		name   string
		prober Prober
		rt     http.RoundTripper
		res    Resource
		err    error
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
			rt: &roundTripper{
				err: &url.Error{},
			},
		},
		{
			name:   "failure: unexpected device",
			prober: &prober{resource: &resource{unexpected: true}},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
		},
		{
			name:   "failure: parsing error",
			prober: &prober{resource: &resource{}},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("}")),
				},
			},
		},
		{
			name:   "success",
			prober: &prober{resource: &resource{}},
			rt: &roundTripper{
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
				Transport: test.rt,
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

func TestTapper_Scan(t *testing.T) {
	tests := []struct {
		name   string
		prober Prober
		rt     http.RoundTripper
		col    Collection
		err    error
	}{
		{
			name: "failure: probe error",
			prober: &prober{
				funcError: &url.Error{},
			},
			err: Errors{},
		},
		{
			name:   "success: empty collection",
			prober: &prober{},
			col:    Collection{},
		},
		{
			name:   "success: collection with one resource",
			prober: &prober{resource: &resource{}},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
			col: Collection{&resource{}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tap := &Tapper{
				probers:   []Prober{test.prober},
				transport: test.rt,
			}

			col, err := tap.Scan([]net.IP{net.ParseIP("192.168.146.123")})

			if !reflect.DeepEqual(col, test.col) {
				t.Fatalf("expected %#v, got %#v", test.col, col)
			}

			var scanError Errors
			switch {
			case errors.As(test.err, &scanError):
				var se *Errors
				if errors.As(err, &se) {
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

func TestTapper_Execute(t *testing.T) {
	tests := []struct {
		name     string
		col      Collection
		proc     procedure
		affected int
		err      error
	}{
		{
			name:     "success: empty collection",
			col:      Collection{},
			affected: 0,
		},
		{
			name: "success: excluded device",
			col:  Collection{&resource{}},
			proc: func(tap *Tapper, res Resource, ch chan<- *ProcedureResult) {
				ch <- &ProcedureResult{
					err: ErrPolicyExcluded,
				}
			},
			affected: 0,
		},
		{
			name: "failure: procedure not supported",
			col:  Collection{&resource{}},
			proc: func(tap *Tapper, res Resource, ch chan<- *ProcedureResult) {
				ch <- &ProcedureResult{
					err: ErrUnsupportedProcedure,
				}
			},
			affected: 0,
			err:      Errors{},
		},
		{
			name: "success",
			col:  Collection{&resource{}},
			proc: func(tap *Tapper, res Resource, ch chan<- *ProcedureResult) {
				ch <- &ProcedureResult{
					dev: res,
				}
			},
			affected: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tap := &Tapper{}

			affected, err := tap.Execute(test.proc, test.col)

			if affected != test.affected {
				t.Fatalf("expected %d affected devices, got %d", test.affected, affected)
			}

			var execError Errors
			switch {
			case errors.As(test.err, &execError):
				var ee Errors
				if errors.As(err, &ee) {
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
