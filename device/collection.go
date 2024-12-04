package device

import (
	"bytes"
	"fmt"
	"sort"
)

// Fields to sort a Collection by
const (
	FieldDriver = "driver"
	FieldIP     = "ip"
	FieldMAC    = "mac"
	FieldName   = "name"
	FieldModel  = "model"
)

// Collection is a slice of device resources.
type Collection []Resource

// Empty returns true if the Collection is empty, false otherwise.
func (c Collection) Empty() bool {
	return len(c) == 0
}

// SortBy a resource field name.
func (c Collection) SortBy(field string) error {
	switch field {
	case FieldDriver:
		sort.Slice(c, func(i, j int) bool {
			return c[i].Driver() < c[j].Driver()
		})

	case FieldIP:
		sort.Slice(c, func(i, j int) bool {
			// Use To16() to ensure all IP addresses are represented in a 16-byte format.
			// This allows consistent comparison regardless of IPv4 or IPv6.
			return bytes.Compare(c[i].IP().To16(), c[j].IP().To16()) < 0
		})

	case FieldMAC:
		sort.Slice(c, func(i, j int) bool {
			return bytes.Compare(c[i].MAC(), c[j].MAC()) < 0
		})

	case FieldName:
		sort.Slice(c, func(i, j int) bool {
			return c[i].Name() < c[j].Name()
		})

	case FieldModel:
		sort.Slice(c, func(i, j int) bool {
			return c[i].Model() < c[j].Model()
		})

	default:
		return fmt.Errorf("%w: %s", ErrInvalidSortByField, field)
	}

	return nil
}
