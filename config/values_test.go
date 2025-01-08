package config

import (
	"encoding/json"
	"errors"
	"github.com/quetzyg/IoTap/device"
	"io"
	"reflect"
	"strings"
	"testing"
)

// badReader returns an error when read.
type badReader struct{}

// Read implements the io.Reader interface.
func (badReader) Read(_ []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestNewValues(t *testing.T) {
	tests := []struct {
		name string
		r    io.Reader
		err  error
	}{
		{
			name: "failure: reader error",
			r:    &badReader{},
			err:  io.ErrUnexpectedEOF,
		},
		{
			name: "failure: syntax error",
			r:    strings.NewReader("}"),
			err:  &json.SyntaxError{},
		},
		{
			name: "success",
			r:    strings.NewReader(`{}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewValues(test.r)

			var syntaxError *json.SyntaxError
			switch {
			case errors.As(test.err, &syntaxError):
				var se *json.SyntaxError
				if errors.As(err, &se) {
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

func TestLoadFromEnv(t *testing.T) {
	tests := []struct {
		name string
		envs map[string]string
		err  error
	}{
		{
			name: "failure: missing credentials",
			err:  errCredentialsNotFound,
		},
		{
			name: "success",
			envs: map[string]string{
				iotapUsername: "foo",
				iotapPassword: "bar",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.envs {
				t.Setenv(k, v)
			}

			_, err := LoadFromEnv()

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}

func TestLoadFromConfigDir(t *testing.T) {
	tests := []struct {
		name string
		dir  string
		err  error
	}{
		{
			name: "no available config directory",
			dir:  "",
		},
		{
			name: "failure: error loading empty file",
			dir:  "invalid/path",
			err:  &json.SyntaxError{},
		},
		{
			name: "success",
			dir:  "../testdata/",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("XDG_CONFIG_HOME", test.dir)
			t.Setenv("HOME", test.dir)

			_, err := LoadFromConfigDir()

			var syntaxError *json.SyntaxError
			switch {
			case errors.As(test.err, &syntaxError):
				var se *json.SyntaxError
				if errors.As(err, &se) {
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

func TestLoadValues(t *testing.T) {
	tests := []struct {
		name string
		envs map[string]string
		val  *Values
		err  error
	}{
		{
			name: "load from env",
			envs: map[string]string{
				iotapUsername: "foo",
				iotapPassword: "bar",
			},
			val: &Values{
				Credentials: &device.Credentials{
					Username: "foo",
					Password: "bar",
				},
			},
		},
		{
			name: "load from file",
			envs: map[string]string{
				"XDG_CONFIG_HOME": "../testdata/",
				"HOME":            "../testdata/",
			},
			val: &Values{
				Credentials: &device.Credentials{
					Username: "admin",
					Password: "secret",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.envs {
				t.Setenv(k, v)
			}

			val, _ := LoadValues()

			if !reflect.DeepEqual(val.Credentials, test.val.Credentials) {
				t.Fatalf("expected %#v, got %#v", test.val.Credentials, val.Credentials)
			}
		})
	}
}
