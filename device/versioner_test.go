package device

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type versioner struct {
	funcError error
	resource
}

func (v *versioner) VersionRequest() (*http.Request, error) {
	if v.funcError != nil {
		return nil, v.funcError
	}

	return http.NewRequest(http.MethodGet, "", nil)
}

func (v *versioner) Outdated() bool {
	return false
}
func (v *versioner) UpdateDetails() string {
	return ""
}

type versionChallenger struct {
	versioner
}

func (vc *versionChallenger) ChallengeAccepted(*http.Response) bool {
	return true
}

func (vc *versionChallenger) ChallengeResponse(r *http.Request, _ *http.Response) (*http.Request, error) {
	return r, nil
}

func TestVersion(t *testing.T) {
	tests := []struct {
		rt   http.RoundTripper
		dev  Resource
		err  error
		name string
	}{
		{
			name: "failure: unsupported procedure",
			dev:  &resource{},
			err:  ErrUnsupportedProcedure,
		},
		{
			name: "failure: function error",
			dev: &versioner{
				funcError: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "failure: http response error",
			dev:  &versioner{},
			rt: &roundTripper{
				err: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "success",
			dev:  &versioner{},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
		},
		{
			name: "success: challenger implementation",
			dev:  &versionChallenger{},
			rt: &roundTripper{
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
				transport: test.rt,
			}

			ch := make(chan *ProcedureResult, 1)

			Version(tap, test.dev, ch)

			result := <-ch

			var urlError *url.Error
			switch {
			case errors.As(test.err, &urlError):
				var ue *url.Error
				if errors.As(result.err, &ue) {
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
