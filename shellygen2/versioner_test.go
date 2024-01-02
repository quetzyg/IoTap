package shellygen2

import (
	"io"
	"net"
	"net/http"
	"reflect"
	"testing"
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

	expectedBody := `{"id":0,"src":"IoTune","method":"Shelly.CheckForUpdate"}`
	if string(body) != expectedBody {
		t.Fatalf("expected %s, got %s", expectedBody, body)
	}
}

func TestDevice_UpdateAvailable(t *testing.T) {
	tests := []struct {
		name      string
		dev       *Device
		available bool
	}{
		{
			name: "update unavailable",
			dev: &Device{
				Version:     "1.0",
				VersionNext: "1.0",
			},
		},
		{
			name: "update available",
			dev: &Device{
				Version:     "1.0",
				VersionNext: "2.0",
			},
			available: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			available := test.dev.UpdateAvailable()
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
				Name:        "stowify-A001",
				Version:     "1.0",
				VersionNext: "2.0",
			},
			details: "[shelly_gen2] stowify-A001 @ 192.168.146.123 can be updated from 1.0 to 2.0",
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
