package config

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
)

// badReadCloser returns an error when read.
type badReadCloser struct{}

// Read implements the io.Reader interface.
func (badReadCloser) Read(_ []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

// Close implements the io.Closer interface.
func (badReadCloser) Close() error {
	return nil
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name string
		r    io.ReadCloser
		val  *Values
		err  error
	}{
		{
			name: "failure: reader error",
			r:    &badReadCloser{},
			val:  &Values{},
			err:  io.ErrUnexpectedEOF,
		},
		{
			name: "failure: syntax error",
			r:    io.NopCloser(strings.NewReader("}")),
			val:  &Values{},
			err:  &json.SyntaxError{},
		},
		{
			name: "success",
			r:    io.NopCloser(strings.NewReader(`{}`)),
			val:  &Values{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := load(test.r, test.val)

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

func TestLoadFromPath(t *testing.T) {
	tests := []struct {
		name string
		fp   string
		val  *Values
		err  error
	}{
		{
			name: "success: invalid file path isn't an error",
			fp:   "/invalid/file/path/is.ok",
		},
		{
			name: "failure: error loading empty file",
			fp:   "../testdata/empty.js",
			err:  &json.SyntaxError{},
		},
		{
			name: "success",
			fp:   "../testdata/values.json",
			val:  &Values{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := LoadFromPath(test.fp, test.val)

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
