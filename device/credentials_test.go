package device

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"strings"
	"testing"
)

func TestNewAuthConfig(t *testing.T) {
	tests := []struct {
		name string
		r    io.Reader
		err  error
	}{
		{
			name: "failure: invalid data",
			r:    strings.NewReader(`}`),
			err:  &json.SyntaxError{},
		},
		{
			name: "success: valid data",
			r:    strings.NewReader(`{"username":"foo","password":"bar"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewAuthConfig(test.r)

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
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := LoadAuthConfig(test.fp)

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
