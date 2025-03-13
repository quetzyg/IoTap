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

func TestDevice_versionUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		dev  *Device
		data []byte
		err  error
	}{
		{
			name: "failure: syntax error",
			dev:  &Device{},
			data: []byte(`}`),
			err:  &json.SyntaxError{},
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

			err := shelly2.versionUnmarshal(test.data)

			var syntaxError *json.SyntaxError
			switch {
			case errors.As(test.err, &syntaxError):
				var se *json.SyntaxError
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

func TestDevice_probeUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		dev  *Device
		data []byte
		err  error
	}{
		{
			name: "failure: syntax error",
			dev:  &Device{},
			data: []byte(`}`),
			err:  &json.SyntaxError{},
		},
		{
			name: "failure: unexpected data",
			dev:  &Device{},
			data: []byte(`{"foo":"bar"}`),
			err:  device.ErrUnexpected,
		},
		{
			name: "failure: bad MAC address",
			dev:  &Device{},
			data: []byte(`{"name":"Shelly Pro 1","id":"shellypro1-001122334455","mac":"foo","model":"SPSW-201XE16EU","gen":2,"fw_id":"20230913-112003/v1.14.0-gcb84623","ver":"1.4.4","app":"Pro1","auth_en":false}`),
			err:  &net.AddrError{},
		},
		{
			name: "success",
			dev:  &Device{},
			data: []byte(`{"name":"Shelly Pro 1","id":"shellypro1-001122334455","mac":"001122334455","model":"SPSW-201XE16EU","gen":2,"fw_id":"20230913-112003/v1.14.0-gcb84623","ver":"1.4.4","app":"Pro1","auth_en":false}`),
		},
		{
			name: "success: device with empty name",
			dev:  &Device{},
			data: []byte(`{"name":null,"id":"shellypro1-001122334455","mac":"001122334455","model":"SPSW-201XE16EU","gen":2,"fw_id":"20230913-112003/v1.14.0-gcb84623","ver":"1.4.4","app":"Pro1","auth_en":false}`),
		},
	}

	shelly2 := &Device{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			err := shelly2.probeUnmarshal(test.data)

			var (
				syntaxError *json.SyntaxError
				addrError   *net.AddrError
			)
			switch {
			case errors.As(test.err, &syntaxError):
				var se *json.SyntaxError
				if errors.As(err, &se) {
					return
				}

			case errors.As(test.err, &addrError):
				var ae *net.AddrError
				if errors.As(err, &ae) {
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

func TestDevice_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		dev  *Device
		data []byte
		err  error
	}{
		{
			name: "success: version unmarshal",
			dev:  &Device{},
			data: []byte(`{"src":"shellypro1-001122334455","result":{"stable":{"version":"20241011-114449/1.4.4-g6d2a586"}}}`),
		},
		{
			name: "success: probe unmarshal",
			dev:  &Device{},
			data: []byte(`{"name":"Shelly Pro 1","id":"shellypro1-001122334455","mac":"001122334455","model":"SPSW-201XE16EU","gen":2,"fw_id":"20230913-112003/v1.14.0-gcb84623","ver":"1.4.4","app":"Pro1","auth_en":false}`),
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
