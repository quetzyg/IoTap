package device

import (
	"bytes"
	"errors"
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

// ErrUnknownSortByField is returned when an attempt is made to
// sort by a field that is not supported by the SortBy() method
var ErrUnknownSortByField = errors.New("unknown field to sort by")

// Collection is a slice of device resources.
type Collection []Resource

// SortBy a resource field name.
func (c Collection) SortBy(field string) ([]Resource, error) {
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
		return nil, fmt.Errorf("%w: %s", ErrUnknownSortByField, field)
	}

	return c, nil
}
