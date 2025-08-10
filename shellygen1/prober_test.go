package shellygen1

import (
	"errors"
	"net"
	"reflect"
	"testing"
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, dev, err := (&Prober{}).Request(test.ip)

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
