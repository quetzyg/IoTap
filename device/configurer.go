package device

import (
	"fmt"
	"net/http"

	"github.com/Stowify/IoTune/httpclient"
)

// Configurer is an interface that provides a standard way to configure IoT devices.
type Configurer interface {
	ConfigureRequests(Config) ([]*http.Request, error)
}

// Configure is a procedure implementation designed to apply configuration settings to an IoT device.
var Configure = func(tun *Tuner, res Resource, ch chan<- *ProcedureResult) {
	dev, ok := res.(Configurer)
	if !ok {
		ch <- &ProcedureResult{
			dev: res,
			err: fmt.Errorf("%w: configure", ErrUnsupportedProcedure),
		}
		return
	}

	rs, err := dev.ConfigureRequests(tun.config)
	if err != nil {
		ch <- &ProcedureResult{
			dev: res,
			err: err,
		}
		return
	}

	client := &http.Client{}

	for _, r := range rs {
		if err = httpclient.Dispatch(client, r, nil); err != nil {
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
