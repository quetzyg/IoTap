package shellygen1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"testing"

	"github.com/quetzyg/IoTap/device"
)

func TestDevice_IP(t *testing.T) {
	ip := net.ParseIP("192.168.146.123")
	shelly1 := &Device{
		ip: ip,
	}

	if !shelly1.IP().Equal(ip) {
		t.Fatalf("expected %v, got %v", ip, shelly1.IP())
	}
}

func TestDevice_MAC(t *testing.T) {
	mac := net.HardwareAddr{00, 17, 34, 51, 68, 85}
	shelly1 := &Device{
		mac: mac,
	}

	if !bytes.Equal(shelly1.MAC(), mac) {
		t.Fatalf("expected %q, got %q", mac, shelly1.MAC())
	}
}

func TestDevice_Name(t *testing.T) {
	name := "Shelly 1"
	shelly1 := &Device{
		name: name,
	}

	if shelly1.Name() != name {
		t.Fatalf("expected %q, got %q", name, shelly1.Name())
	}
}

func TestDevice_Model(t *testing.T) {
	model := "SHSW-1"
	shelly1 := &Device{
		model: model,
	}

	if shelly1.Model() != model {
		t.Fatalf("expected %q, got %q", model, shelly1.Model())
	}
}

func TestDevice_ID(t *testing.T) {
	mac := net.HardwareAddr{00, 17, 34, 51, 68, 85}
	shelly1 := &Device{
		mac: mac,
	}

	if shelly1.ID() != mac.String() {
		t.Fatalf("expected %q, got %q", mac.String(), shelly1.ID())
	}
}

func TestDevice_Driver(t *testing.T) {
	shelly1 := &Device{}

	if shelly1.Driver() != Driver {
		t.Fatalf("expected %q, got %q", Driver, shelly1.Driver())
	}
}

func TestDevice_Vendor(t *testing.T) {
	shelly1 := &Device{}

	if shelly1.Vendor() != Vendor {
		t.Fatalf("expected %q, got %q", Vendor, shelly1.Vendor())
	}
}

func TestDevice_Generation(t *testing.T) {
	shelly1 := &Device{}

	if shelly1.Generation() != Generation {
		t.Fatalf("expected %q, got %q", Generation, shelly1.Generation())
	}
}

func TestDevice_Secured(t *testing.T) {
	shelly1 := &Device{
		secured: true,
	}

	if shelly1.Secured() != true {
		t.Fatalf("expected %t, got %t", true, shelly1.Secured())
	}
}

func TestDevice_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		dev  *Device
		data []byte
		err  error
	}{
		{
			name: "success: version device logic",
			dev:  &Device{},
			data: []byte(`{"new_version":"20230913-112003/v1.14.0-gcb84623"}`),
		},
		{
			name: "failure: unexpected IoT device",
			dev:  &Device{},
			data: []byte(`{}`),
			err:  device.ErrUnexpected,
		},
		{
			name: "success: hydrate device logic",
			dev:  &Device{},
			data: []byte(`{"type":"SHSW-1","mac":"001122334455","auth":false,"fw":"20230913-112003/v1.14.0-gcb84623"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := json.Unmarshal(test.data, test.dev)

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}
