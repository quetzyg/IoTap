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

func TestRegisterConfig(t *testing.T) {
	if len(configRegistry) != 0 {
		t.Fatal("Config registry should be empty")
	}

	RegisterConfig("foo", func() Config {
		return &config{}
	})

	if len(configRegistry) != 1 {
		t.Fatalf("Config registry should have one registered config, %d found", len(configRegistry))
	}
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		r    io.Reader
		cfg  Config
		err  error
		name string
	}{
		{
			name: "failure: syntax error",
			r:    strings.NewReader(`!`),
			err:  &json.SyntaxError{},
		},
		{
			name: "failure: configuration empty",
			r:    strings.NewReader(`{}`),
			err:  ErrConfigurationEmpty,
		},
		{
			name: "success",
			r:    strings.NewReader(`{"foo":"bar"}`),
			cfg: &config{
				Foo: "bar",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := NewConfig(test.r, func() Config {
				return &config{}
			})

			if !reflect.DeepEqual(cfg, test.cfg) {
				t.Fatalf("expected %#v, got %#v", test.cfg, cfg)
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

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		cfg    Config
		err    error
		name   string
		driver string
		fp     string
	}{
		{
			name: "failure: unsupported driver",
			err:  ErrUnsupportedDriver,
		},
		{
			name:   "failure: empty file path",
			driver: "foo",
			err:    ErrFilePathEmpty,
		},
		{
			name:   "failure: file path not found",
			driver: "foo",
			fp:     "foo.bar",
			err:    &fs.PathError{},
		},
		{
			name:   "success",
			driver: "foo",
			fp:     "../testdata/config.json",
			cfg: &config{
				Foo: "bar",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Cleanup(func() {
				configRegistry = make(map[string]ConfigProvider)
			})

			configRegistry["foo"] = func() Config {
				return &config{}
			}

			cfg, err := LoadConfig(test.driver, test.fp)

			if !reflect.DeepEqual(cfg, test.cfg) {
				t.Fatalf("expected %#v, got %#v", test.cfg, cfg)
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
