package device

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"reflect"
	"strings"
	"testing"
)

func TestNewAuthConfig(t *testing.T) {
	tests := []struct {
		name string
		r    io.Reader
		auth *AuthConfig
		err  error
	}{
		{
			name: "failure: invalid data",
			r:    strings.NewReader(`}`),
			err:  &json.SyntaxError{},
		},
		{
			name: "failure: wrong JSON structure",
			r:    strings.NewReader(`{"username":"foo","password":"bar"}`),
			err:  ErrMissingCredentials,
		},
		{
			name: "failure: only username",
			r:    strings.NewReader(`{"credentials":{"username":"foo"}}`),
			err:  ErrMissingCredentials,
		},
		{
			name: "success: only password",
			r:    strings.NewReader(`{"credentials":{"password":"bar"}}`),
			auth: &AuthConfig{
				Policy:      nil,
				Credentials: &Credentials{Username: "", Password: "bar"},
			},
		},
		{
			name: "success: username and password",
			r:    strings.NewReader(`{"credentials":{"username":"foo","password":"bar"}}`),
			auth: &AuthConfig{
				Policy:      nil,
				Credentials: &Credentials{Username: "foo", Password: "bar"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			auth, err := NewAuthConfig(test.r)

			if !reflect.DeepEqual(auth, test.auth) {
				t.Fatalf("expected %#v, got %#v", test.auth, auth)
			}

			var syntaxError *json.SyntaxError
			switch {
			case errors.As(test.err, &syntaxError):
				var se *json.SyntaxError
				if errors.As(err, &se) {
					return
				}

			default:
				if errors.Is(err, test.err) {
					return
				}
			}

			t.Fatalf("expected %#v, got %#v", test.err, err)
		})
	}
}

func TestLoadAuthConfig(t *testing.T) {
	tests := []struct {
		name string
		fp   string
		auth *AuthConfig
		err  error
	}{
		{
			name: "failure: empty file path",
			fp:   "",
			err:  ErrFilePathEmpty,
		},
		{
			name: "failure: file path not found",
			fp:   "foo.bar",
			err:  &fs.PathError{},
		},
		{
			name: "success",
			fp:   "../testdata/authconfig.json",
			auth: &AuthConfig{
				Policy:      nil,
				Credentials: &Credentials{Username: "admin", Password: "secret"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			auth, err := LoadAuthConfig(test.fp)

			if !reflect.DeepEqual(auth, test.auth) {
				t.Fatalf("expected %#v, got %#v", test.auth, auth)
			}

			var pathError *fs.PathError
			switch {
			case errors.As(test.err, &pathError):
				var pe *fs.PathError
				if errors.As(err, &pe) {
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
