package device

import (
	"net"
	"net/http"
)

// Prober defines a method to generate an HTTP request to probe for IoT
// devices on a given IP address. This request is dispatched during the
// network scanning phase.
type Prober interface {
	ProbeRequest(ip net.IP) (*http.Request, Resource, error)
}
