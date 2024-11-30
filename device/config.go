package device

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// Config defines the methods an IoT device configuration instance should implement.
type Config interface {
	Driver() string
	Empty() bool
}

// loadConfig from an I/O reader and unmarshal the data into a Config instance.
func loadConfig(reader io.Reader, config Config) error {
	b, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		return err
	}

	if config.Empty() {
		return ErrConfigurationEmpty
	}

	return nil
}

// LoadConfigFromPath reads and loads a configuration from the specified file path.
// It opens the file, ensures it's closed after reading, and processes the configuration.
// Returns an error if the file cannot be opened or the configuration cannot be loaded.
func LoadConfigFromPath(path string, config Config) error {
	if path == "" {
		return fmt.Errorf("the configuration file path cannot be empty")
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("Config close error: %v", err)
		}
	}()

	return loadConfig(f, config)
}
