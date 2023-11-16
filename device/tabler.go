package device

// Tabler defines the methods for structuring IoT device data
// in a tabular format with consistent column widths and aligned rows.
type Tabler interface {
	ColWidths() [6]int
	Row(format string) string
}
