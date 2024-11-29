package device

import (
	"fmt"
	"net/http"

	"github.com/Stowify/IoTune/httpclient"
)

// Rebooter is an interface that provides a standard way to trigger a reboot on IoT devices.
type Rebooter interface {
	RebootRequest() (*http.Request, error)
}

// Reboot is a procedure implementation designed to reboot an IoT device.
var Reboot = func(_ *Tuner, dev Resource, ch chan<- *ProcedureResult) {
	rsc, ok := dev.(Rebooter)
	if !ok {
		ch <- &ProcedureResult{
			dev: dev,
			err: fmt.Errorf("%w: reboot", ErrUnsupportedProcedure),
		}
		return
	}

	r, err := rsc.RebootRequest()
	if err != nil {
		ch <- &ProcedureResult{
			dev: dev,
			err: err,
		}
		return
	}

	if err = httpclient.Dispatch(&http.Client{}, r, nil); err != nil {
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
