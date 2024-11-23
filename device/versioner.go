package device

import (
	"fmt"
	"net/http"

	"github.com/Stowify/IoTune/httpclient"
)

// Versioner is an interface that provides a set of methods to aid in IoT device versioning.
type Versioner interface {
	VersionRequest() (*http.Request, error)
	UpdateAvailable() bool
	UpdateDetails() string
}

// UpdateDetailsFormat string to be used with fmt.Sprintf()
const UpdateDetailsFormat = "[%s] %s @ %s can be updated from %s to %s"

// Version is a procedure implementation designed to check the version of an IoT device.
var Version = func(tun *Tuner, dev Resource, ch chan<- *ProcedureResult) {
	rsc, ok := dev.(Versioner)
	if !ok {
		ch <- &ProcedureResult{
			dev: dev,
			err: fmt.Errorf("%w: version", ErrUnsupportedProcedure),
		}
		return
	}

	r, err := rsc.VersionRequest()
	if err != nil {
		ch <- &ProcedureResult{
			dev: dev,
			err: err,
		}
		return
	}

	if err = httpclient.Dispatch(&http.Client{}, r, rsc); err != nil {
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
