package device

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

var (
	// ErrConfigurationEmpty is an error type used when loading a configuration file.
	// This error is returned when the configuration file data import results in an empty Config instance.
	ErrConfigurationEmpty = errors.New("empty configuration")

	// ErrUnexpected is an error type used when probing for IoT devices on the network.
	// This error is returned when a device is found, but doesn't match an expected target device.
	ErrUnexpected = errors.New("unexpected IoT device")

	// ErrDriverMismatch is an error type used in the IoT device configuration process.
	// This error is returned when a device tries to use a wrong Config value.
	ErrDriverMismatch = errors.New("device driver mismatch")

	// ErrUnsupportedProcedure is an error type used when executing procedures against IoT devices.
	// This error is returned when a device does not support the procedure that is being executed.
	ErrUnsupportedProcedure = errors.New("unsupported device procedure")

	// ErrStrategyExcluded is an error type used during the configuration strategy logic,
	// excluding a device from being configured.
	ErrStrategyExcluded = errors.New("strategy excluded device")

	// Strategy error types
	errStrategyModeUndefined = errors.New("the strategy mode is undefined")
	errStrategyModeInvalid   = errors.New("the strategy mode is invalid")
)

// ProbeError holds an IP address probe error.
type ProbeError struct {
	ip  net.IP
	err error
}

// Error interface implementation for ProbeError.
func (pe *ProbeError) Error() string {
	return fmt.Sprintf("%s: %v\n", pe.ip, pe.err)
}

// Errors represents a device error collection.
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
