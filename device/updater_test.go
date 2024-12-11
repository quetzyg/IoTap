package device

import (
	"errors"
	"net/http"
	"testing"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		name         string
		roundTripper http.RoundTripper
		dev          Resource
		err          error
	}{
		{
			name: "failure: unsupported procedure",
			dev:  &resource{},
			err:  ErrUnsupportedProcedure,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tap := &Tapper{
				transport: test.roundTripper,
			}

			ch := make(chan *ProcedureResult, 1)

			Update(tap, test.dev, ch)

			result := <-ch

			if !errors.Is(result.err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, result.err)
			}
		})
	}
}
