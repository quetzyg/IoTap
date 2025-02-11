package shellygen2

import (
	"errors"
	"maps"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	"github.com/quetzyg/IoTap/device"
	"github.com/quetzyg/IoTap/httpclient"
)

func TestDevice_ChallengeAccepted(t *testing.T) {
	tests := []struct {
		name     string
		dev      *Device
		resp     *http.Response
		accepted bool
	}{
		{
			name:     "challenge not accepted",
			dev:      &Device{secured: false},
			accepted: false,
		},
		{
			name: "challenge not accepted",
			dev:  &Device{secured: true},
			resp: func() *http.Response {
				resp := &http.Response{
					StatusCode: http.StatusUnauthorized,
					Header:     http.Header{},
				}

				resp.Header.Set(httpclient.WWWAuthenticateHeader, authScheme+" foo=bar")

				return resp
			}(),
			accepted: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			accepted := test.dev.ChallengeAccepted(test.resp)

			if accepted != test.accepted {
				t.Fatalf("expected %t, got %t", test.accepted, accepted)
			}
		})
	}
}

func TestParseDigest(t *testing.T) {
	tests := []struct {
		name string
		resp *http.Response
		dirs map[string]string
		err  error
	}{
		{
			name: "failure: missing digest directives",
			resp: &http.Response{},
			err:  errMissingDigestDirectives,
		},
		{
			name: "success: directives parsed",
			resp: func() *http.Response {
				resp := &http.Response{
					StatusCode: http.StatusUnauthorized,
					Header:     http.Header{},
				}

				resp.Header.Set(
					httpclient.WWWAuthenticateHeader,
					`Digest qop="auth", realm="shellypro1-001122334455", nonce="12345678", algorithm=SHA-256`,
				)

				return resp
			}(),
			dirs: map[string]string{
				"qop":       "auth",
				"realm":     "shellypro1-001122334455",
				"nonce":     "12345678",
				"algorithm": "SHA-256",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dirs, err := parseDigest(test.resp)

			if !maps.Equal(dirs, test.dirs) {
				t.Fatalf("expected %#v, got %#v", test.dirs, dirs)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}

func TestHA1(t *testing.T) {
	const expected = "9753485d35600f865fdc7f84ef1b6f63eea3ee664aa5a5227c7b512ad54d207b"

	hash := ha1("shellypro1-001122334455", "secret")

	if hash != expected {
		t.Fatalf("expected %q, got %q", expected, hash)
	}
}

func TestHA2(t *testing.T) {
	const expected = "44e5ebbe01ba8c658f8af3a9f6ad4f0501d40300e21e69b19ad4669ceb544a14"

	hash := ha2(httptest.NewRequest(http.MethodGet, "/foo/bar", nil))

	if hash != expected {
		t.Fatalf("expected %q, got %q", expected, hash)
	}
}

func compareSecuredRequests(t *testing.T, expected, actual *http.Request) {
	if expected.Method != actual.Method {
		t.Fatalf("expected %q, got %q", expected.Method, actual.Method)
	}

	if expected.URL.String() != actual.URL.String() {
		t.Fatalf("expected %q, got %q", expected.URL.String(), actual.URL.String())
	}

	for header := range expected.Header {
		if actual.Header.Get(header) == "" {
			t.Fatalf("expected header %q to be set", header)
		}

		re := regexp.MustCompile(expected.Header.Get(header))

		if !re.MatchString(actual.Header.Get(header)) {
			t.Fatalf("expected %q to match %q", actual.Header.Get(header), expected.Header.Get(header))
		}
	}
}

func TestDevice_ChallengeResponse(t *testing.T) {
	tests := []struct {
		name string
		dev  *Device
		r1   *http.Request
		resp *http.Response
		r2   *http.Request
		err  error
	}{
		{
			name: "failure: missing credentials",
			dev:  &Device{},
			err:  device.ErrMissingCredentials,
		},
		{
			name: "failure: missing directives",
			dev:  &Device{cred: &device.Credentials{}},
			resp: &http.Response{},
			err:  errMissingDigestDirectives,
		},
		{
			name: "success",
			dev: &Device{ip: net.ParseIP("192.168.146.123"), cred: &device.Credentials{
				Username: "admin",
				Password: "secret",
			}},
			r1: func() *http.Request {
				r := &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   "foo",
					},
					Header: http.Header{},
				}

				return r
			}(),
			resp: func() *http.Response {
				resp := &http.Response{
					StatusCode: http.StatusUnauthorized,
					Header:     http.Header{},
				}

				resp.Header.Set(
					httpclient.WWWAuthenticateHeader,
					`Digest qop="auth", realm="shellypro1-001122334455", nonce="12345678", algorithm=SHA-256`,
				)

				return resp
			}(),
			r2: func() *http.Request {
				r := &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   "foo",
					},
					Header: http.Header{},
				}

				r.Header.Set(
					httpclient.AuthorizationHeader,
					`^Digest username="admin", realm="shellypro1-001122334455", nonce="12345678", uri="foo", response="[a-f0-9]{64}", algorithm=SHA-256, qop=auth, nc=00000001, cnonce="[A-Z2-7]{26}"$`,
				)

				return r
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := test.dev.ChallengeResponse(test.r1, test.resp)

			switch {
			case err == nil:
				compareSecuredRequests(t, test.r2, r)
				return

			case errors.Is(err, test.err):
				return

			default:
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}
