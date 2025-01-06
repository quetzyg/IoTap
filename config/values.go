package config

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/quetzyg/IoTap/device"
)

// File defines the default/expected filename for the configuration.
const File = "iotap.json"

// Values from a loaded configuration file.
type Values struct {
	Credentials *device.Credentials `json:"credentials"`
}

// load data from an I/O reader and unmarshal it into a *Values instance.
func load(r io.Reader, values *Values) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &values)
}

// LoadFromPath reads and loads configuration values from the specified file path.
// It opens the file, ensures it's closed after reading, and processes the configuration data.
// Returns an error if the file cannot be opened or the configuration cannot be loaded.
func LoadFromPath(fp string, values *Values) error {
	f, err := os.Open(fp)
	if err != nil {
		// An unexistent file is not considered an error
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		return err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("Values close error: %v", err)
		}
	}()

	return load(f, values)
}
