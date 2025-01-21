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
		d.Vendor(),
		d.model,
		d.Generation(),
		d.Firmware,
		d.mac.String(),
		fmt.Sprintf("http://%s", d.ip),
		d.name,
		fmt.Sprint(d.secured),
	}, sep)
}

// MarshalJSON implements the Marshaler interface.
func (d *Device) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"vendor":     d.Vendor(),
		"mac":        d.mac.String(),
		"url":        fmt.Sprintf("http://%s", d.ip),
		"name":       d.name,
		"model":      d.model,
		"generation": d.Generation(),
		"firmware":   d.Firmware,
		"secured":    d.secured,
	})
}
