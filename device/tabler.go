package device

// ColWidths is an array holding column width values for six columns.
// Each entry corresponds to the character width of a column, to aid
// in formatting fields for consistent alignment in tabular data representation.
type ColWidths [6]int

// Tabler defines the methods for structuring IoT device data
// in a tabular format with consistent column widths and aligned rows.
type Tabler interface {
	ColWidths() ColWidths
	Row(format string) string
}
