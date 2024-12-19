package device

import (
	"encoding/json"
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

	d.Scripts, err = LoadScriptsFromPath(tmp.Scripts)

	return err
}

// loadDeployment from an I/O reader and unmarshal the data into a *Deployment instance.
func loadDeployment(r io.Reader, dep *Deployment) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &dep)
}

// LoadDeploymentFromPath reads and loads a deployment from the specified file path.
// It opens the file, ensures it's closed after reading, and processes the deployment data.
// Returns an error if the file cannot be opened or the deployment cannot be loaded.
func LoadDeploymentFromPath(fp string, dep *Deployment) error {
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
			log.Printf("Deployment close error: %v", err)
		}
	}()

	return loadDeployment(f, dep)
}
