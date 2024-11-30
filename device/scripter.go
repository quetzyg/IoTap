package device

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Stowify/IoTune/httpclient"
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

// Code returns the content of the script.
func (s *IoTScript) Code() []byte {
	return s.code
}

// Length returns the length of the script content.
func (s *IoTScript) Length() int {
	return len(s.code)
}

// NewIoTScript creates a new *IoTScript instance.
func NewIoTScript(name string) *IoTScript {
	return &IoTScript{
		name: name,
	}
}

// loadScript from an I/O reader and read the data into an *IoTScript instance.
func loadScript(r io.Reader, scr *IoTScript) (err error) {
	scr.code, err = io.ReadAll(r)
	if err != nil {
		return err
	}

	if scr.Length() == 0 {
		return ErrScriptEmpty
	}

	return nil
}

// LoadScriptFromPath reads and loads a script from the specified file path.
// It opens the file, ensures it's closed after reading, and processes the script.
// Returns an error if the file cannot be opened or the script cannot be loaded.
func LoadScriptFromPath(fp string, scr *IoTScript) error {
	if fp == "" {
		return fmt.Errorf("the script file path cannot be empty")
	}

	f, err := os.Open(fp)
	if err != nil {
		return err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("Script close error: %v", err)
		}
	}()

	return loadScript(f, scr)
}

// Scripter is an interface that provides a standard way to set a script on IoT devices.
type Scripter interface {
	ScriptRequests(*IoTScript) ([]*http.Request, error)
}

// Script is a procedure implementation designed to upload a script to an IoT device.
var Script = func(tun *Tuner, dev Resource, ch chan<- *ProcedureResult) {
	rsc, ok := dev.(Scripter)
	if !ok {
		ch <- &ProcedureResult{
			dev: dev,
			err: fmt.Errorf("%w: script", ErrUnsupportedProcedure),
		}
		return
	}

	rs, err := rsc.ScriptRequests(tun.script)
	if err != nil {
		ch <- &ProcedureResult{
			dev: dev,
			err: err,
		}
		return
	}

	client := &http.Client{}

	for _, r := range rs {
		if err = httpclient.Dispatch(client, r, nil); err != nil {
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
