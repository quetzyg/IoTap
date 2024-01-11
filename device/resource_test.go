package device

import (
	"net"
	"testing"
)

func TestMacify(t *testing.T) {
	tests := []struct {
		name     string
		mac      string
		expected string
	}{
		{
			name:     "return invalid mac address format",
			mac:      "ab12cd",
			expected: "ab12cd",
		},
		{
			name:     "return valid mac address format",
			mac:      "ab12cd34ef09",
			expected: "ab:12:cd:34:ef:09",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Macify(test.mac)
			if got != test.expected {
				t.Fatalf("expected %s, got %s", test.expected, got)
			}
		})
	}
}

type devSecured struct{ secure bool }

func (d *devSecured) Driver() string { return "" }

func (d *devSecured) IP() net.IP { return nil }

func (d *devSecured) MAC() net.HardwareAddr { return nil }

func (d *devSecured) ID() string { return "" }

func (d *devSecured) Secured() bool { return d.secure }

func TestSecuredEmoji(t *testing.T) {
	tests := []struct {
		name     string
		dev      *devSecured
		expected string
	}{
		{
			name: "return secured emoji",
			dev: &devSecured{
				secure: true,
			},
			expected: secured,
		},
		{
			name:     "return unsecured emoji",
			dev:      &devSecured{},
			expected: unsecured,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := SecuredEmoji(test.dev)
			if got != test.expected {
				t.Fatalf("expected %s, got %s", test.expected, got)
			}
		})
	}
}
