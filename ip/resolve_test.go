package ip

import (
	"errors"
	"net"
	"testing"
)

func TestValidateNetworkMembership(t *testing.T) {
	tests := []struct {
		name    string
		network *net.IPNet
		err     error
	}{
		{
			name: "failure: nil network",
			err:  errNetworkCannotBeNil,
		},
		{
			name:    "failure: not a network member",
			network: &net.IPNet{},
			err:     errNetworkMembership,
		},
		{
			name: "success",
			network: &net.IPNet{
				IP:   net.ParseIP("127.0.0.1"),
				Mask: net.IPv4Mask(8, 0, 0, 0),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateNetworkMembership(test.network)
			if !errors.Is(err, test.err) {
				t.Fatalf("expected %t, got %t", test.err, err)
			}
		})
	}
}
