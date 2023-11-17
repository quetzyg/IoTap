package device

import (
	"net"
	"net/http"
)

// Prober defines a method to generate an HTTP request for probing an IoT
// device based on a given IP address. This request is dispatched during
// the scanning phase for detecting known devices.
type Prober interface {
	ProbeRequest(ip net.IP) (*http.Request, Resource, error)
}
