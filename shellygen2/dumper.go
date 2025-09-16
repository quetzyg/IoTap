package shellygen2

import (
	"encoding/json/v2"
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
		Vendor:     d.Vendor(),
		Model:      d.model,
		Generation: d.Generation(),
		Firmware:   d.Firmware,
		MAC:        d.mac.String(),
		URL:        fmt.Sprintf("http://%s", d.ip),
		Name:       d.name,
		Secured:    d.secured,
	})
}
