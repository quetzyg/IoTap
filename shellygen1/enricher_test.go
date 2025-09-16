package shellygen1

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"net"
	"net/http"
	"reflect"
	"testing"

	"github.com/quetzyg/IoTap/device"
)

func TestDevice_EnrichRequest(t *testing.T) {
	dev := &Device{
		ip: net.ParseIP("192.168.146.123"),
	}

	r, err := dev.EnrichRequest()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if r.Method != http.MethodGet {
		t.Fatalf("expected %s, got %s", http.MethodGet, r.Method)
	}

	expectedURL := "http://192.168.146.123/settings"
	if r.URL.String() != expectedURL {
		t.Fatalf("expected %s, got %s", expectedURL, r.URL.String())
	}

	expectedHeaders := http.Header{
		"Content-Type": []string{"application/json"},
	}
	if !reflect.DeepEqual(expectedHeaders, r.Header) {
		t.Fatalf("expected %s, got %s", expectedHeaders, r.Header)
	}

	if r.Body != nil {
		t.Fatalf("expected nil, got %#v", r.Body)
	}
}

func TestDevice_EnrichUnmarshaler(t *testing.T) {
	tests := []struct {
		err  error
		dev  *Device
		name string
		data []byte
	}{
		{
			name: "failure: syntactic error",
			dev:  &Device{},
			data: []byte(`}`),
			err:  &jsontext.SyntacticError{},
		},
		{
			name: "failure: missing MAC",
			dev:  &Device{},
			data: []byte(`{"name":"shelly1"}`),
			err:  device.ErrUnexpected,
		},
		{
			name: "success (null name)",
			dev:  &Device{},
			data: []byte(`{"device":{"mac":"001122334455"},"name":null}`),
		},
		{
			name: "success",
			dev: &Device{
				name: "shelly1",
			},
			data: []byte(`{"device":{"mac":"001122334455"},"name":"shelly1"}`),
		},
	}

	shelly1 := &Device{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			err := json.Unmarshal(test.data, shelly1, json.WithUnmarshalers(shelly1.EnrichUnmarshaler()))

			if !reflect.DeepEqual(shelly1, test.dev) {
				t.Fatalf("expected %#v, got %#v", test.dev, shelly1)
			}

			var syntacticError *jsontext.SyntacticError
			switch {
			case errors.As(test.err, &syntacticError):
				var se *jsontext.SyntacticError
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
