package device

import (
	"encoding/json"
	"errors"
	"net"
	"reflect"
	"strings"
	"testing"
)

// resource implementation for testing purposes.
type resource struct {
	driver     string
	ip         net.IP
	mac        net.HardwareAddr
	name       string
	model      string
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
	return "id"
}

// Driver name of this resource implementation.
func (r *resource) Driver() string {
	return r.driver
}

// Secured returns true if the item requires authentication to be accessed, false otherwise.
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
		r.Driver(),
		r.mac.String(),
		r.ip.String(),
		r.Name(),
		r.Model(),
		"v1.2.3",
		SecuredEmoji(r),
	}, sep)
}

// MarshalJSON implements the Marshaler interface.
func (r *resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"driver":   r.Driver(),
		"mac":      r.mac.String(),
		"ip":       r.ip,
		"name":     r.name,
		"model":    r.model,
		"secured":  r.secured,
		"firmware": "v1.2.3",
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
		name  string
		col   Collection
		ord   Collection
		field string
		err   error
	}{
		{
			name: "success: collection sorted by driver",
			col: Collection{
				&resource{driver: "shelly_gen2"},
				&resource{driver: "shelly_gen1"},
			},
			ord: Collection{
				&resource{driver: "shelly_gen1"},
				&resource{driver: "shelly_gen2"},
			},
			field: FieldDriver,
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
