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

// NewAuthConfig creates a new *AuthConfig instance by parsing data from the provided reader.
// It returns an error if the data is invalid or cannot be parsed.
func NewAuthConfig(r io.Reader) (*AuthConfig, error) {
	var auth *AuthConfig
	if err := json.NewDecoder(r).Decode(&auth); err != nil {
		return nil, err
	}

	return auth, nil
}

// LoadAuthConfig creates a new *AuthConfig instance from a file at the given path.
// It returns an error if the file cannot be opened or contains invalid data.
func LoadAuthConfig(fp string) (*AuthConfig, error) {
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
			log.Printf("Auth configuration close error: %v", err)
		}
	}()

	return NewAuthConfig(f)
}
