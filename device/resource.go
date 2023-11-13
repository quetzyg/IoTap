package device

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
)

// Prober defines the methods an implementation should have.
type Prober interface {
	ProbeRequest(ip net.IP) (*http.Request, Resource, error)
}

// Rebooter is an interface that provides a standard way to interact with IoT devices to trigger a reboot.
type Rebooter interface {
	RebootRequest() (*http.Request, error)
}

// Resource defines the methods an IoT device resource should implement.
type Resource interface {
	Driver() string
	IP() net.IP
	ID() string
	UpdateRequest() (*http.Request, error)
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
