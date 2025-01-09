package shellygen2

import (
	"errors"
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"

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

				resp.Header.Set(httpclient.WWWAuthenticateHeader, `Digest qop="auth", realm="shellypro1-001122334455", nonce="12345678", algorithm=SHA-256`)

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

func TestCliNonce(t *testing.T) {
	cnonce := cliNonce()

	if len(cnonce) != 32 {
		t.Fatalf("expected a 32 character length string, got %d", len(cnonce))
	}
}

func TestHA2(t *testing.T) {
	const expected = "44e5ebbe01ba8c658f8af3a9f6ad4f0501d40300e21e69b19ad4669ceb544a14"

	hash := ha2(httptest.NewRequest(http.MethodGet, "/foo/bar", nil))

	if hash != expected {
		t.Fatalf("expected %q, got %q", expected, hash)
	}
}
