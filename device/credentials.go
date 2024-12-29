package device

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// Credentials to interact with secured IoT devices.
type Credentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

// AuthConfig to apply when enabling authentication to IoT devices.
type AuthConfig struct {
	Policy *Policy `json:"policy,omitempty"`
	Credentials
}

// loadAuthConfig from an I/O reader and unmarshal the data into an *AuthConfig instance.
func loadAuthConfig(r io.Reader, cfg *AuthConfig) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &cfg)
}

// LoadAuthConfigFromPath reads and loads an auth configuration from the specified file path.
// It opens the file, ensures it's closed after reading, and processes the auth configuration data.
// Returns an error if the file cannot be opened or the auth configuration cannot be loaded.
func LoadAuthConfigFromPath(fp string, cfg *AuthConfig) error {
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
			log.Printf("Auth configuration close error: %v", err)
		}
	}()

	return loadAuthConfig(f, cfg)
}
