package shellygen2

import (
	"io"
	"net"
	"net/http"
	"reflect"
	"testing"
)

func TestDevice_UpdateRequest(t *testing.T) {
	dev := &Device{
		ip: net.ParseIP("192.168.146.123"),
	}

	r, err := dev.UpdateRequest()
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

	expectedBody := `{"params":{"stage":"stable"},"src":"IoTap","method":"Shelly.Update","id":0}`
	if string(body) != expectedBody {
		t.Fatalf("expected %s, got %s", expectedBody, body)
	}
}
