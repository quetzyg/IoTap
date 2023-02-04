package iot

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
)

// Prober implements methods to help find IoT devices on a network.
type Prober interface {
	MakeRequest(ip net.IP) (*http.Request, Device, error)
	IgnoreError(err error) bool
}

// Device represents an IoT device.
type Device interface {
	Driver() string
	IP() net.IP
	ID() string
}

type Devices map[string]Device

const (
	userAgentHeader = "User-Agent"
	userAgent       = "IoTune/0.1"

	ContentTypeHeader = "Content-Type"
	JSONMimeType      = "application/json"
)

// DeviceFetcher performs an HTTP request to fetch a Device.
func DeviceFetcher(client *http.Client, r *http.Request, dev Device) (Device, error) {
	r.Header.Set(userAgentHeader, userAgent)

	response, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	defer func(body io.ReadCloser) {
		err = body.Close()
		if err != nil {
			log.Printf("error closing response body: %v", err)
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
