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

	dispatcher := httpclient.NewDispatcher(&http.Client{
		Transport: tap.transport,
	})

	var opts []httpclient.DispatchOption

	if challenger, ok := res.(httpclient.Challenger); ok {
		opts = append(opts, httpclient.WithChallenger(challenger))
	}

	if err = dispatcher.Dispatch(r, opts...); err != nil {
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
