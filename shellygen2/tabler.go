package shellygen2

import "fmt"

// ColWidths calculates and returns the width of each of the 6 fields in the Device type
// as an array of integers. The purpose of this method is to assist in correctly formatting
// the string representation of Device data for consistent visual alignment.
func (d *Device) ColWidths() [6]int {
	return [6]int{
		len(d.Driver()),
		len(d.MAC),
		len(d.ip.String()),
		len(d.Firmware),
		len(d.Model),
		len(d.Name),
	}
}

// Row returns a string representation of the IoT device's data. Each field of the data
// is structured in columns, adopting a tabular format suitable for display or incorporation
// into a larger table structure.
func (d *Device) Row(format string) string {
	return fmt.Sprintf(
		format,
		d.Driver(),
		d.MAC,
		d.ip,
		d.Firmware,
		d.Model,
		d.Name,
	)
}
