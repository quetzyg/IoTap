package device

import (
	"fmt"
	"net/http"

	"github.com/quetzyg/IoTap/httpclient"
)

// Deployer is an interface that provides a standard way to deploy a script on supported IoT devices.
type Deployer interface {
	DeployRequests(*http.Client, []*Script) ([]*http.Request, error)
}

// Deploy is a procedure implementation designed to deploy a script to an IoT device.
var Deploy = func(tap *Tapper, res Resource, ch chan<- *ProcedureResult) {
	dev, ok := res.(Deployer)
	if !ok {
		ch <- &ProcedureResult{
			dev: res,
			err: fmt.Errorf("%w: deploy", ErrUnsupportedProcedure),
		}
		return
	}

	client := &http.Client{
		Transport: tap.transport,
	}

	rs, err := dev.DeployRequests(client, tap.scripts)
	if err != nil {
		ch <- &ProcedureResult{
			dev: res,
			err: err,
		}
		return
	}

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
