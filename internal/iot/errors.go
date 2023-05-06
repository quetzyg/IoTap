package iot

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

var ErrWrongDevice = errors.New("wrong IoT device")

// ProbeError holds an IP address probe error.
type ProbeError struct {
	ip  net.IP
	err error
}

// Error interface implementation for ProbeError.
func (pe *ProbeError) Error() string {
	return fmt.Sprintf("%s: %v\n", pe.ip, pe.err)
}

// ProbeErrors represents a ProbeError collection.
type ProbeErrors []*ProbeError

// Error interface implementation for ProbeErrors.
func (pe ProbeErrors) Error() string {
	var s strings.Builder
	for _, e := range pe {
		s.WriteString(e.Error())
	}

	return s.String()
}

// Empty checks if the collection has any errors.
func (pe ProbeErrors) Empty() bool {
	return len(pe) == 0
}

// OperationError holds a device operation error.
type OperationError struct {
	dev Device
	err error
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

// OperationErrors represents a OperationError collection.
type OperationErrors []*OperationError

// Error interface implementation for OperationErrors.
func (oe OperationErrors) Error() string {
	var s strings.Builder
	for _, e := range oe {
		s.WriteString(e.Error())
	}

	return s.String()
}

// Empty checks if the collection has any errors.
func (oe OperationErrors) Empty() bool {
	return len(oe) == 0
}
