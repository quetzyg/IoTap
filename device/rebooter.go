package device

import (
	"fmt"
	"net/http"

	"github.com/quetzyg/IoTap/httpclient"
)

// Rebooter is an interface that provides a standard way to trigger a reboot on IoT devices.
type Rebooter interface {
	RebootRequest() (*http.Request, error)
}

// Reboot is a procedure implementation designed to reboot an IoT device.
var Reboot = func(tap *Tapper, res Resource, ch chan<- *ProcedureResult) {
	dev, ok := res.(Rebooter)
	if !ok {
		ch <- &ProcedureResult{
			dev: res,
			err: fmt.Errorf("%w: reboot", ErrUnsupportedProcedure),
		}
		return
	}

	r, err := dev.RebootRequest()
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

	cha, _ := res.(httpclient.Challenger)

	if err = httpclient.Dispatch(client, r, cha, nil); err != nil {
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
