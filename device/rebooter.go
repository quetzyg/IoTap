package device

import (
	"fmt"
	"net/http"

	iotune "github.com/Stowify/IoTune"
)

// Rebooter is an interface that provides a standard way to trigger a reboot on IoT devices.
type Rebooter interface {
	RebootRequest() (*http.Request, error)
}

// Reboot is a procedure implementation designed to reboot an IoT device.
var Reboot = func(_ *Tuner, dev Resource, ch chan<- *ProcedureResult) {
	res, ok := dev.(Rebooter)
	if !ok {
		ch <- &ProcedureResult{
			dev: dev,
			err: fmt.Errorf("%w: reboot", ErrUnsupportedProcedure),
		}
		return
	}

	r, err := res.RebootRequest()
	if err != nil {
		ch <- &ProcedureResult{
			dev: dev,
			err: err,
		}
		return
	}

	if err = iotune.Dispatch(&http.Client{}, r, nil); err != nil {
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
