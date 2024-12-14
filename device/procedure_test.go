package device

import (
	"errors"
	"net"
	"testing"
)

func TestProcedureResult_Error(t *testing.T) {
	tests := []struct {
		name string
		dev  Resource
		err  error
		out  string
	}{
		{
			name: "error without device details",
			err:  errors.New("some error"),
			out:  "some error",
		},
		{
			name: "error with device details",
			dev: &resource{
				driver: "driver",
				ip:     net.ParseIP("192.168.146.123"),
				mac:    net.HardwareAddr{20, 6, 18, 220, 122, 240},
			},
			err: errors.New("some error"),
			out: "[driver] 14:06:12:dc:7a:f0 @ 192.168.146.123: some error\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			perr := &ProcedureResult{
				dev: test.dev,
				err: test.err,
			}

			if test.out != perr.Error() {
				t.Fatalf("expected %q, got %q", test.out, perr.Error())
			}
		})
	}
}
