package device

import (
	"encoding/json/v2"
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

// ConfigProvider is a function type that returns a Config instance.
type ConfigProvider func() Config

var configRegistry = make(map[string]ConfigProvider)

// RegisterConfig registers a ConfigProvider for a specified driver.
func RegisterConfig(driver string, prov ConfigProvider) {
	configRegistry[driver] = prov
}

// NewConfig creates a Config instance by parsing data from the provided reader and provider function.
// It returns an error if the data is invalid or cannot be parsed.
func NewConfig(r io.Reader, factory ConfigProvider) (Config, error) {
	cfg := factory()

	if err := json.UnmarshalRead(r, &cfg); err != nil {
		return nil, err
	}

	if cfg.Empty() {
		return nil, ErrConfigurationEmpty
	}

	return cfg, nil
}

// LoadConfig creates a new Config instance from a driver name and a file at the given path.
// It returns an error if the file cannot be opened or contains invalid data.
func LoadConfig(driver, fp string) (Config, error) {
	factory, ok := configRegistry[driver]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedDriver, driver)
	}

	if fp == "" {
		return nil, ErrFilePathEmpty
	}

	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("Config close error: %v", err)
		}
	}()

	return NewConfig(f, factory)
}
