package device

// Tabular defines the methods for structuring IoT device data
// in a tabular format with consistent column widths and aligned rows.
type Tabular interface {
	ColWidths() [6]int
	Row(format string) string
}
