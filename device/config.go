package device

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// Config defines the methods an IoT device configuration instance should implement.
type Config interface {
	Driver() string
	Empty() bool
}

// loadConfig from an I/O reader and unmarshal the data into a Config implementation.
func loadConfig(r io.Reader, cfg Config) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	if cfg.Empty() {
		return ErrConfigurationEmpty
	}

	return nil
}

// LoadConfigFromPath reads and loads a configuration from the specified file path.
// It opens the file, ensures it's closed after reading, and processes the configuration data.
// Returns an error if the file cannot be opened or the configuration cannot be loaded.
func LoadConfigFromPath(fp string, cfg Config) error {
	if fp == "" {
		return ErrFilePathEmpty
	}

	f, err := os.Open(fp)
	if err != nil {
		return err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("Config close error: %v", err)
		}
	}()

	return loadConfig(f, cfg)
}
