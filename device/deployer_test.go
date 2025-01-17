package device

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type deployer struct {
	resource
	funcError error
}

func (d *deployer) DeployRequests(*http.Client, *Deployment) ([]*http.Request, error) {
	if d.funcError != nil {
		return nil, d.funcError
	}

	return []*http.Request{
		{
			URL:    &url.URL{},
			Method: http.MethodGet,
		},
	}, nil
}

type deployChallenger struct {
	deployer
}

func (dc *deployChallenger) ChallengeAccepted(*http.Response) bool {
	return true
}

func (dc *deployChallenger) ChallengeResponse(r *http.Request, _ *http.Response) (*http.Request, error) {
	return r, nil
}

func TestDeploy(t *testing.T) {
	tests := []struct {
		name string
		rt   http.RoundTripper
		dev  Resource
		err  error
	}{
		{
			name: "failure: unsupported procedure",
			dev:  &resource{},
			err:  ErrUnsupportedProcedure,
		},
		{
			name: "failure: function error",
			dev: &deployer{
				funcError: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "failure: http response error",
			dev:  &deployer{},
			rt: &roundTripper{
				err: &url.Error{},
			},
			err: &url.Error{},
		},
		{
			name: "success",
			dev:  &deployer{},
			rt: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{}")),
				},
			},
		},
		{
			name: "success: challenger implementation",
			dev:  &deployChallenger{},
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

			Deploy(tap, test.dev, ch)

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
