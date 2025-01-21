package device

import (
	"errors"
	"fmt"
	"net"
)

var (
	// ErrFilePathEmpty indicates that a required file path was not provided.
	// This error is typically returned when an operation requires a file path,
	// but an empty string was passed instead.
	ErrFilePathEmpty = errors.New("the file path cannot be empty")

	// ErrConfigurationEmpty indicates that the loaded configuration file
	// resulted in an empty Config instance.
	ErrConfigurationEmpty = errors.New("the configuration is empty")

	// ErrMissingCredentials is returned when an HTTP request cannot be secured
	// due to missing authentication credentials.
	ErrMissingCredentials = errors.New("missing authentication credentials")

	// ErrScriptEmpty indicates that the loaded script file
	// resulted in an empty Script instance.
	ErrScriptEmpty = errors.New("empty IoT script")

	// ErrUnexpected is returned when unmarshaling payloads that do not conform
	// to the anticipated device structure.
	ErrUnexpected = errors.New("unexpected IoT device")

	// ErrUnsupportedDriver is returned when an unsupported or unrecognised driver is used.
	ErrUnsupportedDriver = errors.New("unsupported driver")

	// ErrDriverMismatch is an error type used in the IoT device configuration process.
	// This error is returned when a device tries to use a wrong Config value.
	ErrDriverMismatch = errors.New("device driver mismatch")

	// ErrUnsupportedProcedure is an error type used when executing procedures against IoT devices.
	// This error is returned when a device does not support the procedure that is being executed.
	ErrUnsupportedProcedure = errors.New("unsupported device procedure")

	// ErrPolicyExcluded is an error type used during the configuration policy logic,
	// if a device is to be excluded from being configured.
	ErrPolicyExcluded = errors.New("policy excluded device")

	// ErrInvalidSortByField is returned when an attempt is made to
	// sort by a field that is not supported by the SortBy() method
	ErrInvalidSortByField = errors.New("invalid field to sort by")

	// ErrInvalidDumpFormat is returned when the dump format isn't correct.
	ErrInvalidDumpFormat = errors.New("invalid dump format")

	// Policy error types
	errPolicyModeUndefined = errors.New("the policy mode is undefined")
	errPolicyModeInvalid   = errors.New("the policy mode is invalid")
)

// ProbeError for an IP address.
type ProbeError struct {
	ip  net.IP
	err error
}

// Error interface implementation.
func (pe *ProbeError) Error() string {
	return fmt.Sprintf("%s: %v", pe.ip, pe.err)
}

// Errors represents an error collection.
type Errors []error

// Error interface implementation.
func (e Errors) Error() string {
	return errors.Join(e...).Error()
}
