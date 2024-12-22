package device

import (
	"net"
	"testing"
)

func TestProbeError_Error(t *testing.T) {
	const expected = "192.168.146.123: use of closed network connection"

	err := &ProbeError{
		ip:  net.ParseIP("192.168.146.123"),
		err: net.ErrClosed,
	}

	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

func TestErrors_Error(t *testing.T) {
	const expected = "unexpected IoT device\ndevice driver mismatch"

	err := &Errors{
		ErrUnexpected,
		ErrDriverMismatch,
	}

	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}
