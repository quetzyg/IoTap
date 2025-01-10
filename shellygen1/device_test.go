package shellygen1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"testing"

	"github.com/quetzyg/IoTap/device"
)

var (
	ip       = net.ParseIP("192.168.146.123")
	mac      = net.HardwareAddr{00, 17, 34, 51, 68, 85}
	name     = "Shelly 1"
	model    = "SHSW-1"
	firmware = "20230913-112003/v1.14.0-gcb84623"
	dev      = &Device{
		ip:       ip,
		mac:      mac,
		name:     name,
		model:    model,
		Firmware: firmware,
	}
)

func TestDevice_IP(t *testing.T) {
	if !dev.IP().Equal(ip) {
		t.Fatalf("expected %v, got %v", ip, dev.IP())
	}
}

func TestDevice_MAC(t *testing.T) {
	if !bytes.Equal(dev.MAC(), mac) {
		t.Fatalf("expected %q, got %q", mac, dev.MAC())
	}
}

func TestDevice_Name(t *testing.T) {
	if dev.Name() != name {
		t.Fatalf("expected %q, got %q", name, dev.Name())
	}
}

func TestDevice_Model(t *testing.T) {
	if dev.Model() != model {
		t.Fatalf("expected %q, got %q", model, dev.Model())
	}
}

func TestDevice_ID(t *testing.T) {
	if dev.ID() != mac.String() {
		t.Fatalf("expected %q, got %q", mac.String(), dev.ID())
	}
}

func TestDevice_Driver(t *testing.T) {
	if dev.Driver() != Driver {
		t.Fatalf("expected %q, got %q", Driver, dev.Driver())
	}
}

func TestDevice_Vendor(t *testing.T) {
	if dev.Vendor() != Vendor {
		t.Fatalf("expected %q, got %q", Vendor, dev.Vendor())
	}
}

func TestDevice_Generation(t *testing.T) {
	if dev.Generation() != Generation {
		t.Fatalf("expected %q, got %q", Generation, dev.Generation())
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
