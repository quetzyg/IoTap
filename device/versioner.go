package device

import (
	"fmt"
	"net/http"

	"github.com/Stowify/IoTune/httpclient"
)

// Versioner is an interface that provides a set of methods to aid in IoT device versioning.
type Versioner interface {
	VersionRequest() (*http.Request, error)
	Outdated() bool
	UpdateDetails() string
}

// UpdateDetailsFormat string to be used with fmt.Sprintf()
const UpdateDetailsFormat = "[%s] %s @ %s can be updated from %s to %s"

// Version is a procedure implementation designed to check the version of an IoT device.
var Version = func(tun *Tuner, res Resource, ch chan<- *ProcedureResult) {
	dev, ok := res.(Versioner)
	if !ok {
		ch <- &ProcedureResult{
			dev: res,
			err: fmt.Errorf("%w: version", ErrUnsupportedProcedure),
		}
		return
	}

	r, err := dev.VersionRequest()
	if err != nil {
		ch <- &ProcedureResult{
			dev: res,
			err: err,
		}
		return
	}

	if err = httpclient.Dispatch(&http.Client{}, r, dev); err != nil {
		ch <- &ProcedureResult{
			dev: res,
			err: err,
		}
		return
	}

	ch <- &ProcedureResult{
		dev: res,
	}
}
