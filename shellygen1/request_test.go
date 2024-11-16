package shellygen1

import (
	"net"
	"net/url"
	"reflect"
	"testing"
)

func TestBuildURL(t *testing.T) {
	ip := net.ParseIP("192.168.146.12")

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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uri := buildURL(ip, test.path)

			if uri != test.expected {
				t.Fatalf("expected %s, got %s", test.expected, uri)
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
