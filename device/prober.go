package device

import (
	"net"
	"net/http"
)

// Prober defines a method to generate an HTTP request to probe for IoT
// devices on a given IP address. This request is dispatched during the
// network scanning phase.
type Prober interface {
	Request(ip net.IP) (*http.Request, Resource, error)
}

// ProberProvider is a function type that returns a Prober instance.
type ProberProvider func() Prober

var proberRegistry = make(map[string]ProberProvider)

// RegisterProber registers a ProberProvider for a specified driver.
func RegisterProber(driver string, prov ProberProvider) {
	proberRegistry[driver] = prov
}

// GetProbers returns a slice of Prober implementations based on the specified driver.
// If the driver is "all", it returns a slice containing all registered probers.
// If a specific driver name is provided, it returns a slice with the prober for that driver.
func GetProbers(driver string) []Prober {
	if driver != AllDrivers {
		return []Prober{
			proberRegistry[driver](),
		}
	}

	var probers []Prober

	for _, prov := range proberRegistry {
		probers = append(probers, prov())
	}

	return probers
}
