package shellygen1

import (
	"net"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_Request(t *testing.T) {
	dev := &Device{
		ip: net.ParseIP("192.168.146.123"),
	}

	r, err := dev.Request()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if r.Method != http.MethodGet {
		t.Fatalf("expected %s, got %s", http.MethodGet, r.Method)
	}

	expectedURL := "http://192.168.146.123/ota"
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

func TestDevice_OutOfDate(t *testing.T) {
	tests := []struct {
		name      string
		dev       *Device
		available bool
	}{
		{
			name: "up to date",
			dev: &Device{
				Firmware:     "1.0",
				FirmwareNext: "1.0",
			},
		},
		{
			name: "out of date",
			dev: &Device{
				Firmware:     "1.0",
				FirmwareNext: "2.0",
			},
			available: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			available := test.dev.OutOfDate()
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
				ip:           net.ParseIP("192.168.146.123"),
				name:         "stowify-A001",
				Firmware:     "1.0",
				FirmwareNext: "2.0",
			},
			details: "[shelly_gen1] stowify-A001 @ 192.168.146.123 can be updated from 1.0 to 2.0",
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
