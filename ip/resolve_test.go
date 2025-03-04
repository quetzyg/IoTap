package ip

import (
	"errors"
	"net"
	"reflect"
	"testing"
)

func TestNext(t *testing.T) {
	tests := []struct {
		name string
		cur  net.IP
		next net.IP
	}{
		{
			name: "success #1",
			cur:  net.ParseIP("192.168.0.0"),
			next: net.ParseIP("192.168.0.1"),
		},
		{
			name: "success #2",
			cur:  net.ParseIP("192.168.0.255"),
			next: net.ParseIP("192.168.1.0"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			next(test.cur)

			if test.cur.String() != test.next.String() {
				t.Fatalf("expected %s, got %s", test.next, test.cur)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	tests := []struct {
		name string
		cidr string
		ips  []net.IP
		err  error
	}{
		{
			name: "failure: invalid CIDR",
			cidr: "192",
			err:  &net.ParseError{},
		},
		{
			name: "success: IP string",
			cidr: "192.168.146.123",
			ips: []net.IP{
				net.ParseIP("192.168.146.123"),
			},
		},
		{
			name: "success: CIDR string",
			cidr: "127.0.0.0/31",
			ips: []net.IP{
				net.ParseIP("127.0.0.0"),
				net.ParseIP("127.0.0.1"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ips, err := Resolve(test.cidr)
			if !reflect.DeepEqual(ips, test.ips) {
				t.Fatalf("expected %#v, got %#v", test.ips, ips)
			}

			var parseError *net.ParseError
			switch {
			case errors.As(test.err, &parseError):
				var pe *net.ParseError
				if errors.As(err, &pe) {
					return
				}

			default:
				if errors.Is(err, test.err) {
					return
				}
			}

			t.Fatalf("expected %#v, got %#v", test.err, err)
		})
	}
}
