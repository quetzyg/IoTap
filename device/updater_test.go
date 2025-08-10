package device

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type updater struct {
	funcError error
	resource
}

func (u *updater) UpdateRequest() (*http.Request, error) {
	if u.funcError != nil {
		return nil, u.funcError
	}

	return http.NewRequest(http.MethodGet, "", nil)
}

type updateChallenger struct {
	updater
}

func (uc *updateChallenger) ChallengeAccepted(*http.Response) bool {
	return true
}

func (uc *updateChallenger) ChallengeResponse(r *http.Request, _ *http.Response) (*http.Request, error) {
	return r, nil
}

func TestUpdate(t *testing.T) {
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
			dev: &updater{
				funcError: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "failure: http response error",
			dev:  &updater{},
			rt: &roundTripper{
				err: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "success",
			dev:  &updater{},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
		},
		{
			name: "success: challenger implementation",
			dev:  &updateChallenger{},
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

			Update(tap, test.dev, ch)

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
