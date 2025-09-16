package shellygen2

import (
	"bytes"
	"encoding/json/v2"
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
				Gen:      2,
				secured:  false,
				Firmware: "20241011-114449/1.4.4-g6d2a586",
			},
			sep: ",",
			row: "Shelly,SPSW-201XE16EU,2,20241011-114449/1.4.4-g6d2a586,00:11:22:33:44:55,http://192.168.146.123,Shelly Pro 1,false",
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
		err  error
		dev  *Device
		name string
		json string
	}{
		{
			name: "success",
			dev: &Device{
				ip:       net.ParseIP("192.168.146.123"),
				mac:      net.HardwareAddr{00, 17, 34, 51, 68, 85},
				name:     "Shelly Pro 1",
				model:    "SPSW-201XE16EU",
				Gen:      2,
				secured:  false,
				Firmware: "20241011-114449/1.4.4-g6d2a586",
			},
			json: `{"vendor":"Shelly","model":"SPSW-201XE16EU","generation":"2","firmware":"20241011-114449/1.4.4-g6d2a586","mac":"00:11:22:33:44:55","url":"http://192.168.146.123","name":"Shelly Pro 1","secured":false}`,
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
