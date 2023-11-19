package device

import (
	"fmt"
	"net/http"

	iotune "github.com/Stowify/IoTune"
)

// Versioner is an interface that provides a set of methods to aid in IoT device versioning.
type Versioner interface {
	VersionRequest() (*http.Request, error)
	VersionOutdated() bool
	UpgradeDetails() string
}

// Version is a procedure implementation designed to check the version of an IoT device.
var Version = func(tun *Tuner, dev Resource, ch chan<- *ProcedureResult) {
	res, ok := dev.(Versioner)
	if !ok {
		ch <- &ProcedureResult{
			err: fmt.Errorf("%w: version", ErrUnsupportedProcedure),
			dev: dev,
		}
		return
	}

	r, err := res.VersionRequest()
	if err != nil {
		ch <- &ProcedureResult{
			dev: dev,
			err: err,
		}
		return
	}

	if err = iotune.Dispatch(&http.Client{}, r, res); err != nil {
		ch <- &ProcedureResult{
			dev: dev,
			err: err,
		}
		return
	}

	ch <- &ProcedureResult{
		dev: dev,
	}
}
