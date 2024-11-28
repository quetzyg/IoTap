package device

import (
	"fmt"
	"net/http"

	"github.com/Stowify/IoTune/httpclient"
)

// Versioner is an interface that provides a set of methods to aid in IoT device versioning.
type Versioner interface {
	Request() (*http.Request, error)
	OutOfDate() bool
	UpdateDetails() string
}

// UpdateDetailsFormat string to be used with fmt.Sprintf()
const UpdateDetailsFormat = "[%s] %s @ %s can be updated from %s to %s"

// Version is a procedure implementation designed to check the version of an IoT device.
var Version = func(tun *Tuner, dev Resource, ch chan<- *ProcedureResult) {
	ver, ok := dev.(Versioner)
	if !ok {
		ch <- &ProcedureResult{
			dev: dev,
			err: fmt.Errorf("%w: version", ErrUnsupportedProcedure),
		}
		return
	}

	r, err := ver.Request()
	if err != nil {
		ch <- &ProcedureResult{
			dev: dev,
			err: err,
		}
		return
	}

	if err = httpclient.Dispatch(&http.Client{}, r, ver); err != nil {
		ch <- &ProcedureResult{
			dev: dev,
			err: err,
		}
		return
	}

	ch <- &ProcedureResult{
		dev: dev,
	}
}
