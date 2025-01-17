package device

import (
	"fmt"
	"net/http"

	"github.com/quetzyg/IoTap/httpclient"
)

// Deployer is an interface that provides a standard way to deploy a script on supported IoT devices.
type Deployer interface {
	DeployRequests(*http.Client, *Deployment) ([]*http.Request, error)
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

	rs, err := dev.DeployRequests(client, tap.deployment)
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
