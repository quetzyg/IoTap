package device

import (
	"fmt"
	"net/http"

	iotune "github.com/Stowify/IoTune"
)

// Updater is an interface that provides a standard way to trigger a firmware update on IoT devices.
type Updater interface {
	UpdateRequest() (*http.Request, error)
}

// Update is a procedure implementation designed to update the firmware of an IoT device.
var Update = func(_ *Tuner, dev Resource, ch chan<- *ProcedureResult) {
	res, ok := dev.(Updater)
	if !ok {
		ch <- &ProcedureResult{
			dev: dev,
			err: fmt.Errorf("%w: update", ErrUnsupportedProcedure),
		}
		return
	}

	r, err := res.UpdateRequest()
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
