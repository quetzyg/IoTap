package device

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"path"
	"reflect"
	"strings"
	"testing"
)

func TestScript(t *testing.T) {
	const (
		fp     = "/foo/bar.js"
		code   = `var foo = "abc";`
		length = len(code)
	)

	s := &Script{
		path: fp,
		code: []byte(code),
	}

	if s.Name() != path.Base(fp) {
		t.Fatalf("expected %q, got %q", path.Base(fp), s.Name())
	}

	if !bytes.Equal(s.Code(), []byte(code)) {
		t.Fatalf("expected %q, got %q", code, s.Code())
	}

	if s.Length() != length {
		t.Fatalf("expected %d, got %d", length, s.Length())
	}
}

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

func TestLoadScript(t *testing.T) {
	tests := []struct {
		name string
		r    io.ReadCloser
		src  *Script
		err  error
	}{
		{
			name: "failure: reader error",
			r:    &badReadCloser{},
			src:  &Script{},
			err:  io.ErrUnexpectedEOF,
		},
		{
			name: "failure: script empty",
			r:    io.NopCloser(strings.NewReader("")),
			src:  &Script{},
			err:  ErrScriptEmpty,
		},
		{
			name: "success",
			r:    io.NopCloser(strings.NewReader(`var foo = "abc";`)),
			src:  &Script{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := loadScript(test.r, test.src)

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}

func TestLoadScriptsFromPath(t *testing.T) {
	tests := []struct {
		name string
		fps  []string
		src  []*Script
		err  error
	}{
		{
			name: "failure: no file paths given #1",
			fps:  nil,
			err:  ErrFilePathEmpty,
		},
		{
			name: "failure: no file paths given #2",
			fps:  []string{},
			err:  ErrFilePathEmpty,
		},
		{
			name: "failure: no file paths given #3",
			fps:  []string{""},
			err:  ErrFilePathEmpty,
		},
		{
			name: "failure: file path not found",
			fps:  []string{"foo.bar"},
			err:  &fs.PathError{},
		},
		{
			name: "failure: error loading empty script",
			fps:  []string{"../testdata/empty.js"},
			err:  ErrScriptEmpty,
		},

		{
			name: "success",
			fps:  []string{"../testdata/script1.js", "../testdata/script2.js"},
			src: []*Script{
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scripts, err := LoadScriptsFromPath(test.fps)

			if !reflect.DeepEqual(scripts, test.src) {
				t.Fatalf("expected %#v, got %#v", test.src, scripts)
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
