package device

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

var (
	// ErrUnexpected is an error type used in the process of probing IoT devices on the network.
	// This error is returned when a device is found on the network but does not match an expected or target device.
	ErrUnexpected = errors.New("unexpected IoT device")

	// ErrUnsupportedProcedure means that a particular procedure is not supported on the IoT device.
	ErrUnsupportedProcedure = errors.New("unsupported device procedure")
)

// ProbeError holds an IP address probe error.
type ProbeError struct {
	ip  net.IP
	err error
}

func NewProbeError(ip net.IP, err error) *ProbeError {
	return &ProbeError{
		ip:  ip,
		err: err,
	}
}

// Error interface implementation for ProbeError.
func (pe *ProbeError) Error() string {
	return fmt.Sprintf("%s: %v\n", pe.ip, pe.err)
}

// OperationError holds a device operation error.
type OperationError struct {
	dev Resource
	err error
}

func NewOperationError(dev Resource, err error) *OperationError {
	return &OperationError{
		dev: dev,
		err: err,
	}
}

// Error interface implementation for OperationError.
func (oe OperationError) Error() string {
	return fmt.Sprintf(
		"[%s] %s @ %s: %v\n",
		oe.dev.Driver(),
		oe.dev.ID(),
		oe.dev.IP(),
		oe.err,
	)
}

// Errors represents an error collection.
type Errors []error

// Error interface implementation for Errors.
func (e Errors) Error() string {
	var s strings.Builder
	for _, err := range e {
		s.WriteString(err.Error())
	}

	return s.String()
}

// Empty checks if the collection has any errors.
func (e Errors) Empty() bool {
	return len(e) == 0
}
