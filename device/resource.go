package device

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

// Prober defines the methods an implementation should have.
type Prober interface {
	ProbeRequest(ip net.IP) (*http.Request, Resource, error)
}

// Resource defines the methods an IoT device resource should implement.
type Resource interface {
	Driver() string
	IP() net.IP
	ID() string
}

// Collection represents a collection of device resources.
type Collection map[string]Resource

// Fetcher performs an HTTP request to fetch a device resource.
func Fetcher(client *http.Client, r *http.Request, dev Resource) (Resource, error) {
	response, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	defer func(body io.ReadCloser) {
		err = body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(response.Body)

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, dev)
	if err != nil {
		return nil, err
	}

	return dev, nil
}

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
