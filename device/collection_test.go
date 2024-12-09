package device

import (
	"net"
	"testing"
)

// resource implementation for testing purposes.
type resource struct {
	secured    bool
	unexpected bool
}

// IP address of the resource.
func (r *resource) IP() net.IP {
	return net.IP{}
}

// MAC address of the resource.
func (r *resource) MAC() net.HardwareAddr {
	return net.HardwareAddr{}
}

// Name of the resource.
func (r *resource) Name() string {
	return "item"
}

// Model of the resource.
func (r *resource) Model() string {
	return "model"
}

// ID returns the resource's unique identifier.
func (r *resource) ID() string {
	return "id"
}

// Driver name of this resource implementation.
func (r *resource) Driver() string {
	return "test"
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
