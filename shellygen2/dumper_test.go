package shellygen2

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"testing"
)

func TestDevice_DelimitedRow(t *testing.T) {
	tests := []struct {
		name string
		dev  *Device
		sep  string
		row  string
	}{
		{
			name: "success",
			dev: &Device{
				ip:       net.ParseIP("192.168.146.123"),
				mac:      net.HardwareAddr{00, 17, 34, 51, 68, 85},
				name:     "Shelly Pro 1",
				model:    "SPSW-201XE16EU",
				secured:  false,
				Firmware: "20241011-114449/1.4.4-g6d2a586",
			},
			sep: ",",
			row: "shelly_gen2,00:11:22:33:44:55,192.168.146.123,Shelly Pro 1,SPSW-201XE16EU,20241011-114449/1.4.4-g6d2a586,false",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			row := test.dev.DelimitedRow(test.sep)

			if row != test.row {
				t.Fatalf("expected %s, got %s", test.row, row)
			}
		})
	}
}

func TestDevice_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		dev  *Device
		json string
		err  error
	}{
		{
			name: "success",
			dev: &Device{
				ip:       net.ParseIP("192.168.146.123"),
				mac:      net.HardwareAddr{00, 17, 34, 51, 68, 85},
				name:     "Shelly Pro 1",
				model:    "SPSW-201XE16EU",
				secured:  false,
				Firmware: "20241011-114449/1.4.4-g6d2a586",
			},
			json: `{"driver":"shelly_gen2","firmware":"20241011-114449/1.4.4-g6d2a586","ip":"192.168.146.123","mac":"00:11:22:33:44:55","model":"SPSW-201XE16EU","name":"Shelly Pro 1","secured":false}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := json.Marshal(test.dev)

			if !bytes.Equal(b, []byte(test.json)) {
				t.Fatalf("expected %q, got %q", test.json, b)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}
		})
	}
}
