package shellygen1

import (
	"github.com/quetzyg/IoTap/device"
	"testing"
)

func TestDevice_Secured(t *testing.T) {
	dev.secured = true

	if dev.Secured() != true {
		t.Fatalf("expected %t, got %t", true, dev.Secured())
	}
}

func TestDevice_SetCredentials(t *testing.T) {
	if dev.cred != nil {
		t.Fatalf("expected nil, got %v", dev.cred)
	}

	dev.SetCredentials(&device.Credentials{
		Username: "admin",
		Password: "admin",
	})

	if dev.cred.Username != "admin" {
		t.Fatalf("expected admin, got %s", dev.cred.Username)
	}

	if dev.cred.Password != "admin" {
		t.Fatalf("expected admin, got %s", dev.cred.Password)
	}
}
