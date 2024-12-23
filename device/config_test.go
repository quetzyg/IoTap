package device

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"strings"
	"testing"
)

// config implementation for testing purposes.
type config struct {
	Foo string `json:"foo"`
}

// Driver name of this Config implementation.
func (c *config) Driver() string {
	return "test"
}

// Empty checks if the struct holding the configuration has a zero value.
func (c *config) Empty() bool {
	return *c == config{}
}

// badReader returns an error when read.
type badReader struct{}

// Read implements the io.Reader interface.
func (badReader) Read(_ []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name string
		r    io.Reader
		cfg  Config
		err  error
	}{
		{
			name: "failure: reader error",
			r:    &badReader{},
			err:  io.ErrUnexpectedEOF,
		},
		{
			name: "failure: syntax error",
			r:    strings.NewReader(`!`),
			err:  &json.SyntaxError{},
		},
		{
			name: "failure: configuration empty",
			r:    strings.NewReader(`{}`),
			cfg:  &config{},
			err:  ErrConfigurationEmpty,
		},
		{
			name: "success: ",
			r:    strings.NewReader(`{"foo":"bar"}`),
			cfg:  &config{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := loadConfig(test.r, test.cfg)

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

func TestLoadConfigFromPath(t *testing.T) {
	tests := []struct {
		name string
		fp   string
		cfg  Config
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
			fp:   "../testdata/config.json",
			cfg:  &config{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := LoadConfigFromPath(test.fp, test.cfg)

			if test.err == nil {
				if err != nil {
					t.Fatalf("expected nil, got %v", err)
				}

				return
			}

			if err == nil {
				t.Fatalf("expected an error but got nil")
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
