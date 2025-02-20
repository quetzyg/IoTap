package shellygen2

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/quetzyg/IoTap/device"
	"github.com/quetzyg/IoTap/httpclient"
)

// Expected authentication authScheme
const authScheme = "Digest"

var errMissingDigestDirectives = errors.New("missing digest directives")

// ChallengeAccepted determines whether the current implementation can accept
// and handle the authentication challenge presented in the provided HTTP response.
func (d *Device) ChallengeAccepted(resp *http.Response) bool {
	if !d.Secured() {
		return false
	}

	return resp.StatusCode == http.StatusUnauthorized &&
		strings.HasPrefix(resp.Header.Get(httpclient.WWWAuthenticateHeader), authScheme)
}

// parseDigest directives from the WWW-Authenticate response header.
func parseDigest(resp *http.Response) (map[string]string, error) {
	scheme, directives, _ := strings.Cut(resp.Header.Get(httpclient.WWWAuthenticateHeader), " ")
	if scheme != authScheme {
		return nil, errMissingDigestDirectives
	}

	dirs := map[string]string{}

	for _, dir := range strings.Split(directives, ", ") {
		k, v, _ := strings.Cut(dir, "=")
		dirs[k] = strings.Trim(v, `"`)
	}

	return dirs, nil
}

// ha1 computes SHA256(username:realm:password) which forms the
// credentials portion of the digest access authentication.
func ha1(realm, password string) string {
	a1 := sha256.Sum256([]byte("admin:" + realm + ":" + password))

	return hex.EncodeToString(a1[:])
}

// ha2 computes SHA256(method:URI) which forms the
// request portion of the digest access authentication.
func ha2(r *http.Request) string {
	a2 := sha256.Sum256([]byte(r.Method + ":" + r.URL.RequestURI()))

	return hex.EncodeToString(a2[:])
}

// ChallengeResponse processes the authentication challenge in the provided
// response and applies the necessary authentication headers to the request.
// See: https://shelly-api-docs.shelly.cloud/gen2/General/Authentication/#authentication-process
func (d *Device) ChallengeResponse(r *http.Request, resp *http.Response) (*http.Request, error) {
	if d.cred == nil {
		return nil, device.ErrMissingCredentials
	}

	dir, err := parseDigest(resp)
	if err != nil {
		return nil, err
	}

	// `nc` is always "00000001", since the `nonce` and `cnonce`
	// are unique per request, making the count unnecessary.
	const nc = "00000001"
	cnonce := rand.Text()

	response := sha256.Sum256([]byte(fmt.Sprintf(
		"%s:%s:%s:%s:%s:%s",
		ha1(dir["realm"], d.cred.Password),
		dir["nonce"],
		nc,
		cnonce,
		dir["qop"],
		ha2(r),
	)))

	r.Header.Set(httpclient.AuthorizationHeader, fmt.Sprintf(
		`Digest username="admin", realm="%s", nonce="%s", uri="%s", response="%s", algorithm=%s, qop=%s, nc=%s, cnonce="%s"`,
		dir["realm"],
		dir["nonce"],
		r.URL.RequestURI(),
		hex.EncodeToString(response[:]),
		dir["algorithm"],
		dir["qop"],
		nc,
		cnonce,
	))

	return r, nil
}
