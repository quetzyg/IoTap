package device

import (
	"fmt"
	"net/http"

	iotune "github.com/Stowify/IoTune"
)

// Configurer is an interface that provides a standard way to configure IoT devices.
type Configurer interface {
	ConfigureRequests(Config) ([]*http.Request, error)
}

// Configure is a procedure implementation designed to apply configuration settings to an IoT device.
var Configure = func(tun *Tuner, dev Resource, ch chan<- *ProcedureResult) {
	rsc, ok := dev.(Configurer)
	if !ok {
		ch <- &ProcedureResult{
			dev: dev,
			err: fmt.Errorf("%w: configure", ErrUnsupportedProcedure),
		}
		return
	}

	rs, err := rsc.ConfigureRequests(tun.config)
	if err != nil {
		ch <- &ProcedureResult{
			dev: dev,
			err: err,
		}
		return
	}

	client := &http.Client{}

	for _, r := range rs {
		if err = iotune.Dispatch(client, r, nil); err != nil {
			ch <- &ProcedureResult{
				dev: dev,
				err: err,
			}
			return
		}
	}

	ch <- &ProcedureResult{
		dev: dev,
	}
}
