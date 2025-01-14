package shellygen1

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/quetzyg/IoTap/device"
	"github.com/quetzyg/IoTap/httpclient"
)

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "success: path without prefix/suffix",
			path:     "foo",
			expected: "http://192.168.146.12/foo",
		},
		{
			name:     "success: path with prefix",
			path:     "/foo",
			expected: "http://192.168.146.12/foo",
		},
		{
			name:     "success: path with suffix",
			path:     "foo/",
			expected: "http://192.168.146.12/foo/",
		},
		{
			name:     "success: path with prefix and suffix",
			path:     "/foo/",
			expected: "http://192.168.146.12/foo/",
		},
	}

	ip := net.ParseIP("192.168.146.12")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uri := buildURL(ip, test.path)

			if uri != test.expected {
				t.Fatalf("expected %s, got %s", test.expected, uri)
			}
		})
	}
}

func TestRequest(t *testing.T) {
	tests := []struct {
		name   string
		dev    *Device
		params url.Values
		r      *http.Request
		err    error
	}{
		{
			name: "success: request with parameters",
			dev:  &Device{ip: net.ParseIP("192.168.146.123")},
			params: url.Values{
				"bar": []string{"baz"},
			},
			r: func() *http.Request {
				r1 := &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme:   "http",
						Host:     "192.168.146.123",
						Path:     probePath,
						RawQuery: "bar=baz",
					},
					Header: http.Header{},
				}

				r1.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

				return r1
			}(),
		},
		{
			name: "failure: credentials missing",
			dev:  &Device{secured: true},
			err:  device.ErrMissingCredentials,
		},
		{
			name: "success: secured request",
			dev: &Device{ip: net.ParseIP("192.168.146.123"), secured: true, cred: &device.Credentials{
				Username: "admin",
				Password: "secret",
			}},
			r: func() *http.Request {
				r1 := &http.Request{
					Method: http.MethodGet,
					URL: &url.URL{
						Scheme: "http",
						Host:   "192.168.146.123",
						Path:   probePath,
					},
					Header: http.Header{},
				}

				r1.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)
				r1.Header.Set(httpclient.AuthorizationHeader, "Basic YWRtaW46c2VjcmV0")

				return r1
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := request(test.dev, probePath, test.params)

			switch {
			case err == nil:
				compareRequests(t, test.r, r)
				return

			case errors.Is(err, test.err):
				return

			default:
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}

func TestSettingsToValues(t *testing.T) {
	tests := []struct {
		name     string
		params   any
		expected url.Values
	}{
		{
			name:     "success: handle unexpected type",
			params:   123,
			expected: nil,
		},
		{
			name: "success: handle nil values",
			params: &settings{
				"foo": nil,
			},
			expected: url.Values{
				"foo": []string{"null"},
			},
		},
		{
			name: "success: handle regular array",
			params: &settings{
				"foo": []any{
					"bar",
					"baz",
					123,
					true,
				},
			},
			expected: url.Values{
				"foo": []string{
					"bar",
					"baz",
					"123",
					"true",
				},
			},
		},
		{
			name: "success: handle special schedule rules array",
			params: &settings{
				"schedule_rules": []any{
					"0800-0123456-on",
					"2000-0123456-off",
				},
			},
			expected: url.Values{
				"schedule_rules": []string{
					"0800-0123456-on,2000-0123456-off",
				},
			},
		},
		{
			name: "success: handle scalar values (string)",
			params: &settings{
				"foo": "bar",
			},
			expected: url.Values{
				"foo": []string{"bar"},
			},
		},
		{
			name: "success: handle scalar values (int)",
			params: &settings{
				"foo": 123,
			},
			expected: url.Values{
				"foo": []string{"123"},
			},
		},
		{
			name: "success: handle scalar values (bool)",
			params: &settings{
				"foo": true,
			},
			expected: url.Values{
				"foo": []string{"true"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			values := settingsToValues(test.params)

			if !reflect.DeepEqual(values, test.expected) {
				t.Fatalf("expected %v, got %v", test.expected, values)
			}
		})
	}
}
