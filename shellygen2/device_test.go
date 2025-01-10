package shellygen2

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
	name     = "Shelly Pro 1"
	model    = "SPSW-201XE16EU"
	firmware = "20241011-114449/1.4.4-g6d2a586"
	dev      = &Device{
		ip:       ip,
		mac:      mac,
		name:     name,
		model:    model,
		Gen:      2,
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
	if dev.Generation() != "2" {
		t.Fatalf("expected %q, got %q", "2", dev.Generation())
	}
}

func TestDevice_Secured(t *testing.T) {
	dev.secured = true

	if dev.Secured() != true {
		t.Fatalf("expected %t, got %t", true, dev.Secured())
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
			data: []byte(`{"src":"shellypro1-001122334455","result":{"stable":{"version":"20241011-114449/1.4.4-g6d2a586"}}}`),
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
			data: []byte(`{"name":"Shelly Pro 1","id":"shellypro1-001122334455","mac":"001122334455","model":"SPSW-201XE16EU","gen":2,"fw_id":"20230913-112003/v1.14.0-gcb84623","ver":"1.4.4","app":"Pro1","auth_en":false}`),
		},
		{
			name: "success: hydrate device logic with empty name",
			dev:  &Device{},
			data: []byte(`{"name":null,"id":"shellypro1-001122334455","mac":"001122334455","model":"SPSW-201XE16EU","gen":2,"fw_id":"20230913-112003/v1.14.0-gcb84623","ver":"1.4.4","app":"Pro1","auth_en":false}`),
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
