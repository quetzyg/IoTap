package device

import (
	"encoding/json/jsontext"
	"errors"
	"io"
	"io/fs"
	"reflect"
	"strings"
	"testing"
)

func TestNewDeployment(t *testing.T) {
	tests := []struct {
		r    io.Reader
		err  error
		name string
	}{
		{
			name: "failure: invalid deployment data",
			r:    strings.NewReader(`}`),
			err:  &jsontext.SyntacticError{},
		},
		{
			name: "success: valid deployment data",
			r:    strings.NewReader(`{"scripts":["../testdata/script1.js"]}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewDeployment(test.r)

			var syntacticError *jsontext.SyntacticError
			switch {
			case errors.As(test.err, &syntacticError):
				var se *jsontext.SyntacticError
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

func TestRegisterDeployer(t *testing.T) {
	if len(deployerRegistry) != 0 {
		t.Fatal("Deployer registry should be empty")
	}

	RegisterDeployer("foo")

	if len(deployerRegistry) != 1 {
		t.Fatalf("Deployer registry should have one registered deployer, %d found", len(configRegistry))
	}
}

func TestLoadDeployment(t *testing.T) {
	tests := []struct {
		err    error
		dep    *Deployment
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
			fp:     "../testdata/deployment.json",
			dep: &Deployment{
				Policy: &Policy{
					Mode:   PolicyModeWhitelist,
					Models: []string{"SHSW-1"},
				},
				Scripts: []*Script{
					{
						path: "../testdata/script1.js",
						code: []byte(`var foo = "abc";`),
					},
					{
						path: "../testdata/script2.js",
						code: []byte(`var bar = 123;`),
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Cleanup(func() {
				deployerRegistry = make(map[string]struct{})
			})

			deployerRegistry["foo"] = struct{}{}

			dep, err := LoadDeployment(test.driver, test.fp)

			if !reflect.DeepEqual(dep, test.dep) {
				t.Fatalf("expected %#v, got %s", test.dep, dep.Scripts[0].code)
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
