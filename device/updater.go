package device

import (
	"fmt"
	"net/http"

	"github.com/quetzyg/IoTap/httpclient"
)

// Updater is an interface that provides a standard way to trigger a firmware update on IoT devices.
type Updater interface {
	UpdateRequest() (*http.Request, error)
}

// Update is a procedure implementation designed to update the firmware of an IoT device.
var Update = func(tap *Tapper, res Resource, ch chan<- *ProcedureResult) {
	dev, ok := res.(Updater)
	if !ok {
		ch <- &ProcedureResult{
			dev: res,
			err: fmt.Errorf("%w: update", ErrUnsupportedProcedure),
		}
		return
	}

	r, err := dev.UpdateRequest()
	if err != nil {
		ch <- &ProcedureResult{
			dev: res,
			err: err,
		}
		return
	}

	client := &http.Client{
		Transport: tap.transport,
	}

	if err = httpclient.Dispatch(client, r, nil); err != nil {
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
