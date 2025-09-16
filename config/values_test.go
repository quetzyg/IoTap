package config

import (
	"encoding/json/jsontext"
	"errors"
	"io"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/quetzyg/IoTap/device"
)

// badReader returns an error when read.
type badReader struct{}

// Read implements the io.Reader interface.
func (badReader) Read(_ []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestNewValues(t *testing.T) {
	tests := []struct {
		r    io.Reader
		err  error
		name string
	}{
		{
			name: "failure: reader error",
			r:    &badReader{},
			err:  io.ErrUnexpectedEOF,
		},
		{
			name: "failure: syntactic error",
			r:    strings.NewReader("}"),
			err:  &jsontext.SyntacticError{},
		},
		{
			name: "success",
			r:    strings.NewReader(`{}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewValues(test.r)

			var syntaxError *jsontext.SyntacticError
			switch {
			case errors.As(test.err, &syntaxError):
				var se *jsontext.SyntacticError
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
		err  error
		envs map[string]string
		name string
	}{
		{
			name: "failure: missing credentials",
			err:  errCredentialsNotFound,
		},
		{
			name: "failure: only username",
			envs: map[string]string{
				iotapUsername: "foo",
			},
			err: errCredentialsNotFound,
		},
		{
			name: "success: only password",
			envs: map[string]string{
				iotapPassword: "bar",
			},
		},
		{
			name: "success: username and password",
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

// Resolve the absolute path for the test $HOME value
var absoluteTestHome = func(t *testing.T) string {
	path, err := filepath.Abs("../testdata/")
	if err != nil {
		t.Fatalf("expected nil, got %#v", err)
	}

	return path
}

func TestLoadFromConfigDir(t *testing.T) {
	tests := []struct {
		err  error
		name string
		dir  string
	}{
		{
			name: "failure: no config directory",
			dir:  "",
			err:  fs.ErrNotExist,
		},
		{
			name: "failure: no config file",
			dir:  "../",
			err:  fs.ErrNotExist,
		},
		{
			name: "success",
			dir:  absoluteTestHome(t),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("XDG_CONFIG_HOME", test.dir)
			t.Setenv("HOME", test.dir)

			_, err := LoadFromConfigDir()

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}

func TestLoadValues(t *testing.T) {
	tests := []struct {
		err  error
		envs map[string]string
		val  *Values
		name string
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
				"XDG_CONFIG_HOME": absoluteTestHome(t),
				"HOME":            absoluteTestHome(t),
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
