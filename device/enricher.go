package device

import (
	"net/http"

	"github.com/quetzyg/IoTap/httpclient"
)

type Enricher interface {
	EnrichRequest() (*http.Request, error)
}

// Enrich is a procedure implementation for device data enrichment.
var Enrich = func(tap *Tapper, res Resource, ch chan<- *ProcedureResult) {
	dev, ok := res.(Enricher)
	if !ok {
		// Device data already complete - no enrichment required
		ch <- &ProcedureResult{
			dev: res,
		}
		return
	}

	r, err := dev.EnrichRequest()
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

	opts := []httpclient.DispatchOption{
		httpclient.WithBinding(dev),
	}

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
