package shellygen2

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
	shelly2 := &Device{
		ip: ip,
	}

	if !shelly2.IP().Equal(ip) {
		t.Fatalf("expected %v, got %v", ip, shelly2.IP())
	}
}

func TestDevice_MAC(t *testing.T) {
	mac := net.HardwareAddr{00, 17, 34, 51, 68, 85}
	shelly2 := &Device{
		mac: mac,
	}

	if !bytes.Equal(shelly2.MAC(), mac) {
		t.Fatalf("expected %q, got %q", mac, shelly2.MAC())
	}
}

func TestDevice_Name(t *testing.T) {
	name := "Shelly Pro 1"
	shelly2 := &Device{
		name: name,
	}

	if shelly2.Name() != name {
		t.Fatalf("expected %q, got %q", name, shelly2.Name())
	}
}

func TestDevice_Model(t *testing.T) {
	model := "SPSW-201XE16EU"
	shelly2 := &Device{
		model: model,
	}

	if shelly2.Model() != model {
		t.Fatalf("expected %q, got %q", model, shelly2.Model())
	}
}

func TestDevice_ID(t *testing.T) {
	mac := net.HardwareAddr{00, 17, 34, 51, 68, 85}
	shelly2 := &Device{
		mac: mac,
	}

	if shelly2.ID() != mac.String() {
		t.Fatalf("expected %q, got %q", mac.String(), shelly2.ID())
	}
}

func TestDevice_Driver(t *testing.T) {
	shelly2 := &Device{}

	if shelly2.Driver() != Driver {
		t.Fatalf("expected %q, got %q", Driver, shelly2.Driver())
	}
}

func TestDevice_Vendor(t *testing.T) {
	shelly2 := &Device{}

	if shelly2.Vendor() != Vendor {
		t.Fatalf("expected %q, got %q", Vendor, shelly2.Vendor())
	}
}

func TestDevice_Generation(t *testing.T) {
	shelly2 := &Device{
		Gen: 2,
	}

	if shelly2.Generation() != "2" {
		t.Fatalf("expected %q, got %q", "2", shelly2.Generation())
	}
}

func TestDevice_Secured(t *testing.T) {
	shelly2 := &Device{
		secured: true,
	}

	if shelly2.Secured() != true {
		t.Fatalf("expected %t, got %t", true, shelly2.Secured())
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
