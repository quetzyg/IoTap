package device

import (
	"net/http"
	"os"
	"path"
)

// IoTScript holds the name and the contents of an IoT device script.
type IoTScript struct {
	name string
	code []byte
}

func (s *IoTScript) Name() string {
	return s.name
}

func (s *IoTScript) Code() []byte {
	return s.code
}

func LoadScript(file string) (*IoTScript, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return &IoTScript{
		name: path.Base(file),
		code: content,
	}, nil
}

// Scripter is an interface that provides a standard way to set a script on IoT devices.
type Scripter interface {
	ScriptRequests(*IoTScript) ([]*http.Request, error)
}
