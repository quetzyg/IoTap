package device

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"strings"
	"testing"
)

func TestLoadDeployment(t *testing.T) {
	tests := []struct {
		name string
		r    io.Reader
		dep  *Deployment
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
			name: "failure: empty script file paths",
			r:    strings.NewReader(`{}`),
			dep:  &Deployment{},
			err:  ErrFilePathEmpty,
		},
		{
			name: "failure: deployment with no scripts",
			r:    strings.NewReader(`{"scripts":["foo.js"]}`),
			dep:  &Deployment{},
			err:  &fs.PathError{},
		},
		{
			name: "success: ",
			r:    strings.NewReader(`{"scripts":["../testdata/script1.js"]}`),
			dep:  &Deployment{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := loadDeployment(test.r, test.dep)

			var (
				syntaxError *json.SyntaxError
				pathError   *fs.PathError
			)
			switch {
			case errors.As(test.err, &syntaxError):
				var se *json.SyntaxError
				if errors.As(err, &se) {
					return
				}

			case errors.As(test.err, &pathError):
				var pe *fs.PathError
				if errors.As(err, &pe) {
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
		dep  *Deployment
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
			dep:  &Deployment{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := LoadDeploymentFromPath(test.fp, test.dep)

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
