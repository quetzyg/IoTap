package shellygen1

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DelimitedRow returns a string representation of the resource,
// with fields separated by a delimiter (e.g., comma or tab).
func (d *Device) DelimitedRow(sep string) string {
	return strings.Join([]string{
		d.Driver(),
		d.mac.String(),
		d.ip.String(),
		d.name,
		d.model,
		d.Firmware,
		fmt.Sprint(d.secured),
	}, sep)
}

// MarshalJSON implements the Marshaler interface.
func (d *Device) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"driver":   d.Driver(),
		"mac":      d.mac.String(),
		"ip":       d.ip,
		"name":     d.name,
		"model":    d.model,
		"firmware": d.Firmware,
		"secured":  d.secured,
	})
}
