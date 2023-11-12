package device

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Config defines the methods an IoT device configuration instance should implement.
type Config interface {
	Driver() string
	MakeRequests(Resource) ([]*http.Request, error)
	Empty() bool
}

var ErrConfigurationEmpty = errors.New("empty configuration")

// LoadConfig from an I/O reader and unmarshal the data into a Config instance.
func LoadConfig(reader io.Reader, config Config) error {
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
