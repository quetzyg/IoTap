package device

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"strings"
	"testing"
)

func TestNewDeployment(t *testing.T) {
	tests := []struct {
		name string
		r    io.Reader
		err  error
	}{
		{
			name: "failure: invalid deployment data",
			r:    strings.NewReader(`}`),
			err:  &json.SyntaxError{},
		},
		{
			name: "success: valid deployment data",
			r:    strings.NewReader(`{"scripts":["../testdata/script1.js"]}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewDeployment(test.r)

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

func TestLoadDeploymentFromPath(t *testing.T) {
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
			fp:   "../testdata/deployment.json",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := LoadDeployment(test.fp)

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
