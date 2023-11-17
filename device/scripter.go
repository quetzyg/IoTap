package device

import (
	"fmt"
	"net/http"
	"os"
	"path"

	iotune "github.com/Stowify/IoTune"
)

// IoTScript holds the name and the contents of an IoT device script.
type IoTScript struct {
	name string
	code []byte
}

// Name of the file the script was loaded from.
func (s *IoTScript) Name() string {
	return s.name
}

// Code returns the contents of the script.
func (s *IoTScript) Code() []byte {
	return s.code
}

// LoadScript from a local file path.
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

// Script is a procedure implementation designed to upload a script to an IoT device.
var Script = func(tun *Tuner, dev Resource, ch chan<- *ProcedureResult) {
	res, ok := dev.(Scripter)
	if !ok {
		ch <- &ProcedureResult{
			err: fmt.Errorf("%w: script", ErrUnsupportedProcedure),
			dev: dev,
		}
		return
	}

	rs, err := res.ScriptRequests(tun.script)
	if err != nil {
		ch <- &ProcedureResult{
			dev: dev,
			err: err,
		}
		return
	}

	client := &http.Client{}

	for _, r := range rs {
		if err = iotune.Dispatch(client, r, nil); err != nil {
			ch <- &ProcedureResult{
				dev: dev,
				err: err,
			}
			return
		}
	}

	ch <- &ProcedureResult{
		dev: dev,
	}
}
