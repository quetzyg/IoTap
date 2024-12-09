package device

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"testing"
)

func TestDumpCSV(t *testing.T) {
	tests := []struct {
		name string
		col  Collection
		out  string
		sep  string
		err  error
	}{
		{
			name: "success: comma separator",
			col: Collection{
				&resource{
					driver:  "shelly_gen1",
					ip:      net.ParseIP("192.168.146.123"),
					mac:     net.HardwareAddr{00, 17, 34, 51, 68, 85},
					name:    "Storage",
					model:   "SHSW-1",
					secured: false,
				},
			},
			out: "shelly_gen1,00:11:22:33:44:55,192.168.146.123,Storage,SHSW-1,v1.2.3,🔓\n",
			sep: ",",
		},
		{
			name: "success: tab separator",
			col: Collection{
				&resource{
					driver:  "shelly_gen1",
					ip:      net.ParseIP("192.168.146.123"),
					mac:     net.HardwareAddr{00, 17, 34, 51, 68, 85},
					name:    "Storage",
					model:   "SHSW-1",
					secured: false,
				},
			},
			out: "shelly_gen1  00:11:22:33:44:55  192.168.146.123  Storage  SHSW-1  v1.2.3  🔓\n",
			sep: "\t",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := &bytes.Buffer{}

			err := dumpCSV(test.col, w, test.sep)
			if !errors.Is(err, test.err) {
				t.Fatalf("expected %v, got %v", test.err, err)
			}

			if !strings.Contains(w.String(), test.out) {
				t.Fatalf("expected %q, got %q", test.out, w.String())
			}
		})
	}
}

func TestDumpJSON(t *testing.T) {
	tests := []struct {
		name string
		col  Collection
		out  string
		err  error
	}{
		{
			name: "success",
			col: Collection{
				&resource{
					driver:  "shelly_gen1",
					ip:      net.ParseIP("192.168.146.123"),
					mac:     net.HardwareAddr{00, 17, 34, 51, 68, 85},
					name:    "Storage",
					model:   "SHSW-1",
					secured: false,
				},
			},
			out: `[
  {
    "driver": "shelly_gen1",
    "firmware": "v1.2.3",
    "ip": "192.168.146.123",
    "mac": "00:11:22:33:44:55",
    "model": "SHSW-1",
    "name": "Storage",
    "secured": false
  }
]`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := &bytes.Buffer{}

			err := dumpJSON(test.col, w)
			if !errors.Is(err, test.err) {
				t.Fatalf("expected %v, got %v", test.err, err)
			}

			if !strings.Contains(w.String(), test.out) {
				t.Fatalf("expected %q, got %q", test.out, w.String())
			}
		})
	}
}
