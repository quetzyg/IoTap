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

// ConfigError holds a device configuration error.
type ConfigError struct {
	dev Device
	err error
}

// Error interface implementation for ConfigError.
func (ce ConfigError) Error() string {
	return fmt.Sprintf(
		"[%s] %s @ %s: %v\n",
		ce.dev.Driver(),
		ce.dev.ID(),
		ce.dev.IP(),
		ce.err,
	)
}

// ConfigErrors represents a ConfigError collection.
type ConfigErrors []*ConfigError

// Error interface implementation for ConfigErrors.
func (ce ConfigErrors) Error() string {
	var s strings.Builder
	for _, e := range ce {
		s.WriteString(e.Error())
	}

	return s.String()
}

// Empty checks if the collection has any errors.
func (ce ConfigErrors) Empty() bool {
	return len(ce) == 0
}
