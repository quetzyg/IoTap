package shellygen1

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
				name:     "Shelly 1",
				model:    "SHSW-1",
				secured:  false,
				Firmware: "20230913-112003/v1.14.0-gcb84623",
			},
			sep: ",",
			row: "Shelly,00:11:22:33:44:55,192.168.146.123,Shelly 1,SHSW-1,1,20230913-112003/v1.14.0-gcb84623,false",
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
				name:     "Shelly 1",
				model:    "SHSW-1",
				secured:  false,
				Firmware: "20230913-112003/v1.14.0-gcb84623",
			},
			json: `{"firmware":"20230913-112003/v1.14.0-gcb84623","generation":"1","ip":"192.168.146.123","mac":"00:11:22:33:44:55","model":"SHSW-1","name":"Shelly 1","secured":false,"vendor":"Shelly"}`,
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
