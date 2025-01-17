package device

import (
	"fmt"
	"net/http"

	"github.com/quetzyg/IoTap/httpclient"
)

// Configurer is an interface that provides a standard way to configure IoT devices.
type Configurer interface {
	ConfigureRequests(Config) ([]*http.Request, error)
}

// Configure is a procedure implementation designed to apply configuration settings to an IoT device.
var Configure = func(tap *Tapper, res Resource, ch chan<- *ProcedureResult) {
	dev, ok := res.(Configurer)
	if !ok {
		ch <- &ProcedureResult{
			dev: res,
			err: fmt.Errorf("%w: configure", ErrUnsupportedProcedure),
		}
		return
	}

	rs, err := dev.ConfigureRequests(tap.config)
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

	for _, r := range rs {
		if err = dispatcher.Dispatch(r, opts...); err != nil {
			ch <- &ProcedureResult{
				dev: res,
				err: err,
			}
			return
		}
	}

	ch <- &ProcedureResult{
		dev: res,
	}
}
