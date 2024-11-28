package device

import (
	"fmt"
	"net/http"

	"github.com/Stowify/IoTune/httpclient"
)

// Updater is an interface that provides a standard way to trigger a firmware update on IoT devices.
type Updater interface {
	UpdateRequest() (*http.Request, error)
}

// Update is a procedure implementation designed to update the firmware of an IoT device.
var Update = func(_ *Tuner, dev Resource, ch chan<- *ProcedureResult) {
	rsc, ok := dev.(Updater)
	if !ok {
		ch <- &ProcedureResult{
			dev: dev,
			err: fmt.Errorf("%w: update", ErrUnsupportedProcedure),
		}
		return
	}

	r, err := rsc.UpdateRequest()
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

// ExecUpdate encapsulates the execution of the device.Update procedure.
func ExecUpdate(tuner *Tuner, devices Collection) error {
	if len(devices) == 0 {
		return nil
	}

	return tuner.Execute(Update)
}
