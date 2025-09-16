package device

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"log"
	"os"
)

// Deployment holds a policy to enforce when deploying scripts to one or more IoT devices.
type Deployment struct {
	Policy  *Policy   `json:"policy,omitempty"`
	Scripts []*Script `json:"scripts"`
}

// UnmarshalJSON implements the Unmarshaler interface.
func (d *Deployment) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Policy  *Policy  `json:"policy,omitempty"`
		Scripts []string `json:"scripts"`
	}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	d.Policy = tmp.Policy
	d.Scripts, err = LoadScripts(tmp.Scripts)

	return err
}

var deployerRegistry = make(map[string]struct{})

// RegisterDeployer registers a ConfigProvider for a specified driver.
func RegisterDeployer(driver string) {
	deployerRegistry[driver] = struct{}{}
}

// NewDeployment creates a new *Deployment instance by parsing data from the provided reader.
// It returns an error if the data is invalid or cannot be parsed.
func NewDeployment(r io.Reader) (*Deployment, error) {
	var dep Deployment
	if err := json.UnmarshalRead(r, &dep); err != nil {
		return nil, err
	}

	return &dep, nil
}

// LoadDeployment creates a new *Deployment instance from a file at the given path.
// It returns an error if the file cannot be opened or contains invalid data.
func LoadDeployment(driver, fp string) (*Deployment, error) {
	if _, ok := deployerRegistry[driver]; !ok {
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
			log.Printf("Deployment close error: %v", err)
		}
	}()

	return NewDeployment(f)
}
