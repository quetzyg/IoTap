package shellygen1

import (
	"net"
	"net/http"
	"reflect"
	"testing"
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
