package device

import (
	"io"
	"net/http"
)

// IoTScript holds the name and the contents of an IoT device script.
type IoTScript struct {
	name string
	code io.Reader
}

func (s *IoTScript) Name() string {
	return s.name
}

func (s *IoTScript) Code() io.Reader {
	return s.code
}

// Scripter is an interface that provides a standard way to set a script on IoT devices.
type Scripter interface {
	ScriptRequests(*IoTScript) ([]*http.Request, error)
}
