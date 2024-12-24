package device

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// Credentials to apply when enabling security on IoT devices.
type Credentials struct {
	Policy   *Policy `json:"policy,omitempty"`
	Username string  `json:"username"`
	Password string  `json:"password"`
}

// loadCredentials from an I/O reader and unmarshal the data into a *Credentials instance.
func loadCredentials(r io.Reader, cred *Credentials) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &cred)
}

// LoadCredentialsFromPath reads and loads credentials from the specified file path.
// It opens the file, ensures it's closed after reading, and processes the credentials data.
// Returns an error if the file cannot be opened or the credentials cannot be loaded.
func LoadCredentialsFromPath(fp string, cred *Credentials) error {
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
			log.Printf("Credentials close error: %v", err)
		}
	}()

	return loadCredentials(f, cred)
}
