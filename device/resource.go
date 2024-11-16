package device

import (
	"net"
	"strings"
)

// Driver is the default value when the user does not specify which driver it wants to target.
const Driver = "all"

// Resource defines the methods an IoT device resource should implement.
type Resource interface {
	Driver() string
	IP() net.IP
	MAC() net.HardwareAddr
	ID() string
	Secured() bool
}

// Collection represents a collection of device resources.
type Collection []Resource

// Macify takes a non-delimited string representation of a MAC address as input
// and returns a properly formatted MAC address with appropriate colon (:) separators.
func Macify(address string) string {
	// If we don't have exactly 12 characters, just return what we got
	if len(address) != 12 {
		return address
	}

	var octets []string
	for i := 0; i < len(address); i += 2 {
		octets = append(octets, address[i:i+2])
	}

	return strings.Join(octets, ":")
}

const (
	secured   = "🔒"
	unsecured = "🔓"
)

// SecuredEmoji returns an open/closed lock emoji, depending on the secured device state.
func SecuredEmoji(dev Resource) string {
	if dev.Secured() {
		return secured
	}

	return unsecured
}
