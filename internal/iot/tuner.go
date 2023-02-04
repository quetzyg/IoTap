package iot

import (
	"errors"
	"net"
	"net/http"
	"net/url"
)

// Tuner holds Devices found during a network scan.
// It also has the ability to set their configuration.
type Tuner struct {
	devices Devices
}

// NewTuner creates a new Tuner instance.
func NewTuner() *Tuner {
	return &Tuner{
		devices: Devices{},
	}
}

// Probe an IP address for known IoT devices.
func Probe(client *http.Client, ip net.IP, prs ...ProbeRequest) (Device, error) {
	for _, pr := range prs {
		r, dev, err := pr(ip)
		if err != nil {
			return nil, err
		}

		dev, err = DeviceFetcher(client, r, dev)
		var ue *url.Error
		if errors.As(err, &ue) {
			// Ignore timeouts, refused connections and other classic HTTP shenanigans,
			// since - normally - it will mean that there's not device at that IP.
			continue
		}

		// The Device might be nil, so the caller will have to validate it
		return dev, err
	}

	return nil, nil
}

// ScanResult represents the outcome of an IP address scan operation.
type ScanResult struct {
	dev Device
	err *ProbeError
}

// scanIP probes a specific IP and passes the result to a channel.
func scanIP(ch chan<- *ScanResult, ip net.IP, prs ...ProbeRequest) {
	result := &ScanResult{}

	dev, err := Probe(&http.Client{}, ip, prs...)
	if err != nil {
		result.err = &ProbeError{ip: ip, err: err}
		ch <- result
		return
	}

	// No device found
	if dev == nil {
		ch <- result
		return
	}

	result.dev = dev
	ch <- result
	return
}

// The usable IP addresses of a /24 subnet.
const subnet24 = 254

// Scan the network for IoT devices using Probe requests.
func (t *Tuner) Scan(ip net.IP, prs ...ProbeRequest) error {
	// Cleanup before scanning
	t.devices = Devices{}

	ch := make(chan *ScanResult)

	var octet byte
	for octet = 1; octet <= subnet24; octet++ {
		go scanIP(ch, net.IPv4(ip[0], ip[1], ip[2], octet), prs...)
	}

	errs := ProbeErrors{}

	for i := 0; i < subnet24; i++ {
		select {
		case result := <-ch:
			if result.err != nil {
				errs = append(errs, result.err)
			}

			if result.dev != nil {
				t.devices[result.dev.ID()] = result.dev
			}
		}
	}

	close(ch)

	return errs
}

// Devices that were found during the network scan.
func (t *Tuner) Devices() Devices {
	return t.devices
}
