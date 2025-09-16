package shellygen2

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"net"
	"reflect"
	"testing"

	"github.com/quetzyg/IoTap/device"
)

func TestProber_Request(t *testing.T) {
	tests := []struct {
		err  error
		dev  *Device
		name string
		uri  string
		ip   net.IP
	}{
		{
			name: "success",
			ip:   net.ParseIP("192.168.146.123"),
			dev: &Device{
				ip: net.ParseIP("192.168.146.123"),
			},
			uri: "http://192.168.146.123/shelly",
			err: nil,
		},
	}

	prober := &Prober{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, dev, err := prober.Request(test.ip)

			if r.URL.String() != test.uri {
				t.Fatalf("expected %#v, got %#v", test.uri, r.URL.String())
			}

			if !reflect.DeepEqual(dev, test.dev) {
				t.Fatalf("expected %#v, got %#v", test.dev, dev)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}

func TestProber_Unmarshaler(t *testing.T) {
	tests := []struct {
		err  error
		dev  *Device
		name string
		data []byte
	}{
		{
			name: "failure: syntactic error",
			dev:  &Device{},
			data: []byte(`}`),
			err:  &jsontext.SyntacticError{},
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			err := json.Unmarshal(test.data, test.dev, json.WithUnmarshalers((&Prober{}).Unmarshaler()))

			var (
				syntacticError *jsontext.SyntacticError
				addrError      *net.AddrError
			)
			switch {
			case errors.As(test.err, &syntacticError):
				var se *jsontext.SyntacticError
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
