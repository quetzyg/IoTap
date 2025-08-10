package httpclient

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
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
func (badReadCloser) Close() error { return nil }

func TestCloneRequest(t *testing.T) {
	tests := []struct {
		err  error
		r    *http.Request
		name string
	}{
		{
			name: "failure: bad body",
			r: &http.Request{
				Body: badReadCloser{},
			},
			err: io.ErrUnexpectedEOF,
		},
		{
			name: "success: all OK",
			r:    httptest.NewRequest(http.MethodGet, "http://localhost/foo?a=b&c=d", strings.NewReader(`{"foo":"bar"}`)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			clone, err := cloneRequest(test.r)

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}

			if err == nil {
				// Compare URL
				if test.r.URL.String() != clone.URL.String() {
					t.Fatalf("expected %#v, got %#v", test.r, clone)
				}

				// Compare HTTP Method
				if test.r.Method != clone.Method {
					t.Fatalf("expected %#v, got %#v", test.r, clone)
				}

				// Compare Headers
				if !reflect.DeepEqual(test.r.Header, clone.Header) {
					t.Fatalf("expected %#v, got %#v", test.r, clone)
				}

				// Compare Body
				body1, _ := io.ReadAll(test.r.Body)
				body2, _ := io.ReadAll(clone.Body)

				if !bytes.Equal(body1, body2) {
					t.Fatalf("expected %#v, got %#v", test.r, clone)
				}
			}
		})
	}
}
