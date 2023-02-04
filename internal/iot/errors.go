package iot

import (
	"fmt"
	"net"
	"strings"
)

// ProbeError holds an IP address probe error.
type ProbeError struct {
	ip  net.IP
	err error
}

// Error interface implementation for ProbeError.
func (pe ProbeError) Error() string {
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
