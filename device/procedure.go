package device

import "fmt"

// procedure is a function type that encapsulates operations to be carried out on IoT devices.
type procedure func(tap *Tapper, res Resource, ch chan<- *ProcedureResult)

// ProcedureResult encapsulates the outcome of a procedure executed on an IoT device.
// These can be related to various operations such as probing, updating, rebooting or configuring a device.
type ProcedureResult struct {
	dev Resource
	err error
}

// Failed checks if the ProcedureResult execution has failed.
func (pr *ProcedureResult) Failed() bool {
	return pr.err != nil
}

// Error interface implementation for ProcedureResult.
func (pr *ProcedureResult) Error() string {
	if pr.dev == nil {
		return pr.err.Error()
	}

	return fmt.Sprintf(
		"[%s] %s @ %s: %v",
		pr.dev.Driver(),
		pr.dev.ID(),
		pr.dev.IP(),
		pr.err,
	)
}
