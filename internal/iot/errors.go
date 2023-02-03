package iot

import (
	"fmt"
	"net"
	"strings"
)

// ScanError holds an IP address scanning error.
type ScanError struct {
	ip  net.IP
	err error
}

// Error implements the error interface for ScanError.
func (se ScanError) Error() string {
	return fmt.Sprintf("%s: %v\n", se.ip, se.err)
}

// ScanErrors represents a ScanError collection.
type ScanErrors []*ScanError

// Error implements the error interface for ScanErrors.
func (se ScanErrors) Error() string {
	var s strings.Builder
	for _, e := range se {
		s.WriteString(e.Error())
	}

	return s.String()
}
