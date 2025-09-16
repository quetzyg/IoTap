package device

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"
	"testing"
)

// resource implementation for testing purposes.
type resource struct {
	vendor     string
	driver     string
	name       string
	model      string
	gen        string
	ip         net.IP
	mac        net.HardwareAddr
	secured    bool
	unexpected bool
}

// IP address of the resource.
func (r *resource) IP() net.IP {
	return r.ip
}

// MAC address of the resource.
func (r *resource) MAC() net.HardwareAddr {
	return r.mac
}

// Name of the resource.
func (r *resource) Name() string {
	return r.name
}

// Model of the resource.
func (r *resource) Model() string {
	return r.model
}

// ID returns the resource's unique identifier.
func (r *resource) ID() string {
	return r.mac.String()
}

// Driver name of this resource implementation.
func (r *resource) Driver() string {
	return r.driver
}

// Vendor represents the name of the company that developed the device.
func (r *resource) Vendor() string {
	return r.vendor
}

// Generation represents the generation of this device.
func (r *resource) Generation() string {
	return r.gen
}

// Secured returns true if the device requires authentication to be accessed, false otherwise.
func (r *resource) Secured() bool {
	return r.secured
}

// UnmarshalJSON implements the Unmarshaler interface.
func (r *resource) UnmarshalJSON(_ []byte) error {
	if r.unexpected {
		return ErrUnexpected
	}

	return nil
}

// DelimitedRow returns a string representation of the resource,
// with fields separated by a delimiter (e.g., comma or tab).
func (r *resource) DelimitedRow(sep string) string {
	return strings.Join([]string{
		r.Vendor(),
		r.model,
		r.Generation(),
		"v1.2.3",
		r.mac.String(),
		fmt.Sprintf("http://%s", r.ip),
		r.name,
		fmt.Sprint(r.secured),
	}, sep)
}

// MarshalJSON implements the Marshaler interface.
func (r *resource) MarshalJSON() ([]byte, error) {
	type v struct {
		Vendor     string `json:"vendor"`
		Model      string `json:"model"`
		Generation string `json:"generation"`
		Firmware   string `json:"firmware"`
		MAC        string `json:"mac"`
		URL        string `json:"url"`
		Name       string `json:"name"`
		Secured    bool   `json:"secured"`
	}

	return json.Marshal(v{
		Vendor:     r.Vendor(),
		Model:      r.model,
		Generation: r.Generation(),
		Firmware:   "v1.2.3",
		MAC:        r.mac.String(),
		URL:        fmt.Sprintf("http://%s", r.ip),
		Name:       r.name,
		Secured:    r.secured,
	})
}

func TestCollection_Empty(t *testing.T) {
	tests := []struct {
		name  string
		col   Collection
		empty bool
	}{
		{
			name:  "collection is empty (nil)",
			col:   nil,
			empty: true,
		},
		{
			name:  "collection is empty",
			col:   Collection{},
			empty: true,
		},
		{
			name: "collection has items",
			col: Collection{
				&resource{},
			},
			empty: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			empty := test.col.Empty()
			if empty != test.empty {
				t.Fatalf("expected %t, got %t", test.empty, empty)
			}
		})
	}
}

func TestCollection_SortBy(t *testing.T) {
	tests := []struct {
		err   error
		name  string
		field string
		col   Collection
		ord   Collection
	}{
		{
			name: "success: collection sorted by vendor",
			col: Collection{
				&resource{vendor: "Shelly"},
				&resource{vendor: "Tuya"},
				&resource{vendor: "SONOFF"},
			},
			ord: Collection{
				&resource{vendor: "SONOFF"},
				&resource{vendor: "Shelly"},
				&resource{vendor: "Tuya"},
			},
			field: FieldVendor,
		},
		{
			name: "success: collection sorted by ip",
			col: Collection{
				&resource{ip: net.ParseIP("192.168.146.200")},
				&resource{ip: net.ParseIP("192.168.146.100")},
			},
			ord: Collection{
				&resource{ip: net.ParseIP("192.168.146.100")},
				&resource{ip: net.ParseIP("192.168.146.200")},
			},
			field: FieldIP,
		},
		{
			name: "success: collection sorted by mac address",
			col: Collection{
				&resource{mac: net.HardwareAddr{102, 119, 136, 153, 170, 187}},
				&resource{mac: net.HardwareAddr{00, 17, 34, 51, 68, 85}},
			},
			ord: Collection{
				&resource{mac: net.HardwareAddr{00, 17, 34, 51, 68, 85}},
				&resource{mac: net.HardwareAddr{102, 119, 136, 153, 170, 187}},
			},
			field: FieldMAC,
		},
		{
			name: "success: collection sorted by name",
			col: Collection{
				&resource{name: "Kitchen"},
				&resource{name: "Office"},
				&resource{name: "Garage"},
			},
			ord: Collection{
				&resource{name: "Garage"},
				&resource{name: "Kitchen"},
				&resource{name: "Office"},
			},
			field: FieldName,
		},
		{
			name: "success: collection sorted by model",
			col: Collection{
				&resource{model: "SPSW-201XE16EU"},
				&resource{model: "SHSW-1"},
				&resource{model: "SNSW-001X16EU"},
			},
			ord: Collection{
				&resource{model: "SHSW-1"},
				&resource{model: "SNSW-001X16EU"},
				&resource{model: "SPSW-201XE16EU"},
			},
			field: FieldModel,
		},
		{
			name: "success: collection sorted by generation",
			col: Collection{
				&resource{gen: "2"},
				&resource{gen: "3"},
				&resource{gen: "1"},
			},
			ord: Collection{
				&resource{gen: "1"},
				&resource{gen: "2"},
				&resource{gen: "3"},
			},
			field: FieldGeneration,
		},
		{
			name:  "failure: invalid sort field",
			field: "foo",
			err:   ErrInvalidSortByField,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.col.SortBy(test.field)
			if !errors.Is(err, test.err) {
				t.Fatalf("expected %#v, got %#v", test.err, err)
			}

			if !reflect.DeepEqual(test.ord, test.col) {
				t.Fatalf("expected %#v, got %#v", test.ord, test.col)
			}
		})
	}
}
