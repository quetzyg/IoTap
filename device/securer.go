package device

import (
	"fmt"
	"net/http"

	"github.com/quetzyg/IoTap/httpclient"
)

// Securer deals with IoT device security.
type Securer interface {
	Secured() bool
	AuthConfigRequest(*AuthConfig) (*http.Request, error)
	SecureRequest(*http.Request) (*http.Request, error)
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

	r, err := dev.AuthConfigRequest(tap.authConfig)
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
