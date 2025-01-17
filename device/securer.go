package device

import (
	"fmt"
	"net/http"

	"github.com/quetzyg/IoTap/httpclient"
)

// Securer deals with IoT device security.
type Securer interface {
	SetCredentials(*Credentials)
	AuthConfigRequest(*AuthConfig) (*http.Request, error)
}

// Secure is a procedure implementation for securing an IoT device.
var Secure = func(tap *Tapper, res Resource, ch chan<- *ProcedureResult) {
	dev, ok := res.(Securer)
	if !ok {
		ch <- &ProcedureResult{
			dev: res,
			err: fmt.Errorf("%w: secure", ErrUnsupportedProcedure),
		}
		return
	}

	r, err := dev.AuthConfigRequest(tap.auth)
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
