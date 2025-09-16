package device

import (
	"encoding/json/v2"
	"fmt"
	"net/http"

	"github.com/quetzyg/IoTap/httpclient"
)

// Versioner is an interface that provides a set of methods to aid in IoT device versioning.
type Versioner interface {
	VersionRequest() (*http.Request, error)
	VersionUnmarshaler() *json.Unmarshalers
	Outdated() bool
	UpdateDetails() string
}

// UpdateDetailsFormat string to be used with fmt.Sprintf()
const UpdateDetailsFormat = "[%s] %s @ %s can be updated from %s to %s"

// Version is a procedure implementation designed to check the version of an IoT device.
var Version = func(tap *Tapper, res Resource, ch chan<- *ProcedureResult) {
	dev, ok := res.(Versioner)
	if !ok {
		ch <- &ProcedureResult{
			dev: res,
			err: fmt.Errorf("%w: version", ErrUnsupportedProcedure),
		}
		return
	}

	r, err := dev.VersionRequest()
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
		httpclient.WithUnmarshaler(dev.VersionUnmarshaler()),
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
