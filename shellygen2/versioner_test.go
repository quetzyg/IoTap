package shellygen2

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"io"
	"net"
	"net/http"
	"reflect"
	"testing"

	"github.com/quetzyg/IoTap/device"
)

func TestDevice_VersionRequest(t *testing.T) {
	dev := &Device{
		ip: net.ParseIP("192.168.146.123"),
	}

	r, err := dev.VersionRequest()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if r.Method != http.MethodPost {
		t.Fatalf("expected %s, got %s", http.MethodPost, r.Method)
	}

	expectedURL := "http://192.168.146.123/rpc"
	if r.URL.String() != expectedURL {
		t.Fatalf("expected %s, got %s", expectedURL, r.URL.String())
	}

	expectedHeaders := http.Header{
		"Content-Type": []string{"application/json"},
	}
	if !reflect.DeepEqual(expectedHeaders, r.Header) {
		t.Fatalf("expected %s, got %s", expectedHeaders, r.Header)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	expectedBody := `{"src":"IoTap","method":"Shelly.CheckForUpdate","id":0}`
	if string(body) != expectedBody {
		t.Fatalf("expected %s, got %s", expectedBody, body)
	}
}

func TestDevice_Outdated(t *testing.T) {
	tests := []struct {
		dev       *Device
		name      string
		available bool
	}{
		{
			name: "up to date",
			dev: &Device{
				Version:     "1.0",
				VersionNext: "1.0",
			},
		},
		{
			name: "out of date",
			dev: &Device{
				Version:     "1.0",
				VersionNext: "2.0",
			},
			available: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			available := test.dev.Outdated()
			if available != test.available {
				t.Fatalf("expected %v, got %v", test.available, available)
			}
		})
	}
}

func TestDevice_UpdateDetails(t *testing.T) {
	tests := []struct {
		name    string
		dev     *Device
		details string
	}{
		{
			name: "no update details",
			dev:  &Device{},
		},
		{
			name: "update details",
			dev: &Device{
				ip:          net.ParseIP("192.168.146.123"),
				name:        "DEV-A001",
				Version:     "1.0",
				VersionNext: "2.0",
			},
			details: "[shellygen2] DEV-A001 @ 192.168.146.123 can be updated from 1.0 to 2.0",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			details := test.dev.UpdateDetails()
			if details != test.details {
				t.Fatalf("expected %s, got %s", test.details, details)
			}
		})
	}
}

func TestDevice_VersionUnmarshaler(t *testing.T) {
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
			name: "failure: unexpected data",
			dev:  &Device{},
			data: []byte(`{"foo":"bar"}`),
			err:  device.ErrUnexpected,
		},
		{
			name: "success: skip beta version",
			dev:  &Device{},
			data: []byte(`{"result":{"beta":{"version":"1.5.1-beta2","build_id":"20250310-083328/1.5.1-beta2-g322cd2a"}}}`),
		},
		{
			name: "success",
			dev:  &Device{},
			data: []byte(`{"result":{"stable":{"version":"20241011-114449/1.4.4-g6d2a586"}}}`),
		},
	}

	shelly2 := &Device{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			err := json.Unmarshal(test.data, shelly2, json.WithUnmarshalers(shelly2.VersionUnmarshaler()))

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
