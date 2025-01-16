package device

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/tabwriter"
)

// Device dump formats
const (
	FormatCSV  = "csv"
	FormatJSON = "json"
)

// Dumper defines an interface for serializing resource data into different output formats.
type Dumper interface {
	DelimitedRow(sep string) string
	MarshalJSON() ([]byte, error)
}

// dumpCSV writes the given Collection of devices to the provided io.Writer in CSV format.
// Each device's data is serialized using the DelimitedRow() method.
func dumpCSV(devices Collection, w io.Writer, sep string) error {
	writer := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	header := []string{
		"Vendor",
		"MAC Address",
		"URL",
		"Name",
		"Model",
		"Generation",
		"Firmware",
		"Secured",
	}

	_, err := fmt.Fprintln(writer, strings.Join(header, sep))
	if err != nil {
		return err
	}

	for _, device := range devices {
		_, err = fmt.Fprintln(writer, device.(Dumper).DelimitedRow(sep))
		if err != nil {
			return err
		}
	}

	defer func() {
		err = writer.Flush()
		if err != nil {
			log.Fatalf("Writer flush error: %v", err)
		}
	}()

	return nil
}

// dumpJSON writes the given Collection of devices to the provided io.Writer in JSON format.
// Each device's data is serialized using the MarshalJSON() method.
func dumpJSON(devices Collection, w io.Writer) error {
	b, err := json.MarshalIndent(devices, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write(b)

	return err
}

// ExecDump is a wrapper function to easily dump device scan results to multiple formats and outputs.
func ExecDump(devices Collection, format string, file string) error {
	var (
		w   io.Writer = os.Stdout
		err error
	)

	if file != "" {
		w, err = os.Create(file)
		if err != nil {
			return err
		}
	}

	switch format {
	case FormatCSV:
		sep := ","

		// Use the Tab separator when outputting to a screen
		if w == os.Stdout {
			sep = "\t"
		}

		return dumpCSV(devices, w, sep)

	case FormatJSON:
		return dumpJSON(devices, w)

	default:
		return fmt.Errorf("%w: %s", ErrInvalidDumpFormat, format)
	}
}
