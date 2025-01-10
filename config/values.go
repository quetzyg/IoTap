package config

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/quetzyg/IoTap/device"
)

const (
	// file defines the default filename for the IoTap configuration.
	file = "iotap.json"

	// ENV variable names
	iotapUsername = "IOTAP_USERNAME"
	iotapPassword = "IOTAP_PASSWORD"
)

// Values from an IoTap configuration.
type Values struct {
	Credentials *device.Credentials `json:"credentials,omitempty"`
}

// NewValues creates a new *Values instance by parsing data from the provided reader.
// It returns an error if the data is invalid or cannot be parsed.
func NewValues(r io.Reader) (*Values, error) {
	var val *Values
	if err := json.NewDecoder(r).Decode(&val); err != nil {
		return nil, err
	}

	return val, nil
}

var errCredentialsNotFound = errors.New("credentials not found in the environment")

// LoadFromEnv initializes a new *Values instance using values from the environment.
// It returns an error if the required environment variables are missing or invalid.
func LoadFromEnv() (*Values, error) {
	username := os.Getenv(iotapUsername)
	password := os.Getenv(iotapPassword)

	// Certain devices only require a password for authentication and do not use a username
	if password == "" {
		return nil, errCredentialsNotFound
	}

	return &Values{
		Credentials: &device.Credentials{
			Username: username,
			Password: password,
		},
	}, nil
}

// LoadFromConfigDir creates a new *Values instance from a user config directory file.
// It returns an error if the file cannot be opened or contains invalid data.
func LoadFromConfigDir() (*Values, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		// Missing config directories are ignored and not treated as errors.
		return nil, nil
	}

	f, err := os.Open(filepath.Join(dir, file))
	if err != nil {
		// Missing files are ignored and not treated as errors.
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}

		return nil, err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("Values close error: %v", err)
		}
	}()

	return NewValues(f)
}

// LoadValues creates a new *Values instance from two sources, with order of precedence:
// 1. Environment variables (IOTAP_*)
// 2. Configuration file at default location (~/.config/iotap.json)
func LoadValues() (*Values, error) {
	val, err := LoadFromEnv()
	if err == nil {
		return val, nil
	}

	return LoadFromConfigDir()
}
