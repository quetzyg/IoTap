package shellygen2

import (
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/quetzyg/IoTap/device"
	"github.com/quetzyg/IoTap/httpclient"
)

// roundTripper is a custom type used for mocking HTTP responses.
type roundTripper struct {
	response *http.Response
	err      error
}

// RoundTrip implements the roundTripper interface.
func (rt roundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return rt.response, rt.err
}

func TestFetchScripts(t *testing.T) {
	tests := []struct {
		name         string
		roundTripper http.RoundTripper
		scripts      []*script
		err          error
	}{
		{
			name: "failure: dispatch failed",
			roundTripper: &roundTripper{
				err: net.ErrClosed,
			},
			err: net.ErrClosed,
		},
		{
			name: "success: scripts retrieved",
			roundTripper: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"result":{"scripts":[{"id":1}]}}`)),
				},
			},
			scripts: []*script{
				{ID: 1},
			},
		},
	}

	shelly2 := &Device{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scripts, err := shelly2.fetchScripts(&http.Client{Transport: test.roundTripper})

			if !reflect.DeepEqual(scripts, test.scripts) {
				t.Fatalf("expected %#v, got %#v", test.scripts, scripts)
			}

			switch {
			case errors.Is(err, test.err):
				return

			default:
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}

func TestDevice_DeployRequests(t *testing.T) {
	tests := []struct {
		name         string
		dep          *device.Deployment
		roundTripper http.RoundTripper
		rs           []*http.Request
		err          error
	}{
		{
			name: "failure: excluded via policy",
			dep: &device.Deployment{
				Policy: &device.Policy{
					Mode:   device.PolicyModeWhitelist,
					Models: []string{"SPSW-201XE16EU"},
				},
			},
			err: device.ErrPolicyExcluded,
		},
		{
			name: "failure: unable to fetch scripts",
			dep:  &device.Deployment{},
			roundTripper: &roundTripper{
				err: net.ErrClosed,
			},
			err: net.ErrClosed,
		},
		{
			name: "failure: unable to fetch scripts",
			dep: &device.Deployment{
				Scripts: func() []*device.Script {
					src, err := device.LoadScript("../testdata/script1.js")
					if err != nil {
						t.Fatalf("unable to load script: %v", err)
					}

					return []*device.Script{
						src,
					}
				}(),
			},
			roundTripper: &roundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"result":{"scripts":[{"id":1}]}}`)),
				},
			},
			rs: func() []*http.Request {
				r1 := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   rpcPath,
					},
					Header: http.Header{},
					Body:   io.NopCloser(bytes.NewBufferString(`{"id":0,"src":"IoTap","method":"Script.Delete","params":{"id":1}}`)),
				}

				r1.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				r2 := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   rpcPath,
					},
					Header: http.Header{},
					Body:   io.NopCloser(bytes.NewBufferString(`{"id":0,"src":"IoTap","method":"Script.Create","params":{"name":"script1.js"}}`)),
				}

				r2.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				r3 := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   rpcPath,
					},
					Header: http.Header{},
					Body:   io.NopCloser(bytes.NewBufferString(`{"id":0,"src":"IoTap","method":"Script.PutCode","params":{"append":false,"code":"var foo = \"abc\";","id":1}}`)),
				}

				r3.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				r4 := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   rpcPath,
					},
					Header: http.Header{},
					Body:   io.NopCloser(bytes.NewBufferString(`{"id":0,"src":"IoTap","method":"Script.SetConfig","params":{"config":{"enable":true},"id":1}}`)),
				}

				r4.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				r5 := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   rpcPath,
					},
					Header: http.Header{},
					Body:   io.NopCloser(bytes.NewBufferString(`{"id":0,"src":"IoTap","method":"Script.Start","params":{"id":1}}`)),
				}

				r5.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				r6 := &http.Request{
					Method: http.MethodPost,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   rpcPath,
					},
					Header: http.Header{},
					Body:   io.NopCloser(bytes.NewBufferString(`{"id":0,"src":"IoTap","method":"Shelly.Reboot"}`)),
				}

				r6.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				return []*http.Request{r1, r2, r3, r4, r5, r6}
			}(),
		},
	}

	shelly2 := &Device{ip: net.ParseIP("192.168.146.123")}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &http.Client{
				Transport: test.roundTripper,
			}
			rs, err := shelly2.DeployRequests(client, test.dep)

			for i, r := range rs {
				compareRequests(t, test.rs[i], r)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}
