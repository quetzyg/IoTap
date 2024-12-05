package device

import (
	"net"
	"testing"
)

// Item is a resource implementation for testing purposes.
type Item struct {
	secured bool
}

// IP address of the Item.
func (i *Item) IP() net.IP {
	return net.IP{}
}

// MAC address of the Item.
func (i *Item) MAC() net.HardwareAddr {
	return net.HardwareAddr{}
}

// Name of the Item.
func (i *Item) Name() string {
	return "item"
}

// Model of the Item.
func (i *Item) Model() string {
	return "model"
}

// ID returns the Item's unique identifier.
func (i *Item) ID() string {
	return "id"
}

// Driver name of this Item implementation.
func (i *Item) Driver() string {
	return "test"
}

// Secured returns true if the item requires authentication to be accessed, false otherwise.
func (i *Item) Secured() bool {
	return i.secured
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
				&Item{},
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
