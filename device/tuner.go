package device

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
)

// Tuner holds Devices found during a network scan.
// It also has the ability to set their configuration.
type Tuner struct {
	devices Collection
}

// NewTuner creates a new tuner instance.
func NewTuner() *Tuner {
	return &Tuner{
		devices: Collection{},
	}
}

// Probe an IP address for a specific IoT device.
func Probe(client *http.Client, ip net.IP, prober Prober) (Resource, error) {
	r, dev, err := prober.ProbeRequest(ip)
	if err != nil {
		return nil, err
	}

	dev, err = Fetcher(client, r, dev)

	var ue *url.Error
	if errors.As(err, &ue) {
		// Ignore timeouts, refused connections and other classic HTTP shenanigans,
		// since (NORMALLY!) it means there's no such device at the IP address.
		return dev, nil
	}

	if errors.Is(err, ErrUnexpected) {
		// Skip unexpected devices.
		return dev, nil
	}

	var je *json.SyntaxError
	if errors.As(err, &je) {
		// We found something, but it's not outputting valid JSON
		return dev, nil
	}

	return dev, err
}

// ProbeResult represents the outcome of an IP address probe operation.
type ProbeResult struct {
	dev Resource
	err *ProbeError
}

// probe probes a specific IP and passes the result to a channel.
func probe(ch chan<- *ProbeResult, ip net.IP, prober Prober) {
	result := &ProbeResult{}

	dev, err := Probe(&http.Client{}, ip, prober)
	if err != nil {
		result.err = NewProbeError(ip, err)
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

// Scan the network with an IoT device prober.
func (t *Tuner) Scan(ip net.IP, prober Prober) error {
	// Cleanup before scanning
	t.devices = Collection{}

	ch := make(chan *ProbeResult)

	for octet := byte(1); octet <= subnet24; octet++ {
		go probe(ch, net.IPv4(ip[0], ip[1], ip[2], octet), prober)
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
func (t *Tuner) Devices() Collection {
	return t.devices
}

// pushConfig handles a single HTTP configuration push to a device.
func pushConfig(client *http.Client, r *http.Request) error {
	response, err := client.Do(r)
	if err != nil {
		return err
	}

	defer func(body io.ReadCloser) {
		err = body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(response.Body)

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		b, _ := io.ReadAll(response.Body)
		return fmt.Errorf("%s - %s (%d)", r.URL.Path, b, response.StatusCode)
	}

	return nil
}

// OperationResult represents the outcome of a device operation.
type OperationResult struct {
	dev      Resource
	finished bool
	err      error
}

// configure a single device.
func configure(ch chan<- *OperationResult, cfg Config, dev Resource) {
	rs, err := cfg.MakeRequests(dev)
	if err != nil {
		ch <- &OperationResult{
			dev:      dev,
			finished: true,
			err:      err,
		}
		return
	}

	client := &http.Client{}

	for _, r := range rs {
		if err = pushConfig(client, r); err != nil {
			ch <- &OperationResult{
				dev:      dev,
				finished: true,
				err:      err,
			}
			return
		}

		ch <- &OperationResult{
			dev: dev,
		}
	}

	ch <- &OperationResult{
		dev:      dev,
		finished: true,
	}
}

// ConfigureDevices found in the network.
func (t *Tuner) ConfigureDevices(cfg Config) error {
	ch := make(chan *OperationResult)

	for _, dev := range t.devices {
		go configure(ch, cfg, dev)
	}

	errs := OperationErrors{}

	remaining := len(t.devices)

	for remaining != 0 {
		select {
		case result := <-ch:
			if result.finished {
				remaining--
			}

			if result.err != nil {
				errs = append(errs, NewOperationError(result.dev, result.err))
			}
		}
	}

	close(ch)

	return errs
}

// update a single device.
func update(ch chan<- *OperationResult, dev Resource) {
	client := &http.Client{}

	r, err := dev.UpdateRequest()
	if err != nil {
		ch <- &OperationResult{
			dev:      dev,
			finished: true,
			err:      err,
		}
	}

	if err = pushConfig(client, r); err != nil {
		ch <- &OperationResult{
			dev:      dev,
			finished: true,
			err:      err,
		}
		return
	}

	ch <- &OperationResult{
		dev:      dev,
		finished: true,
	}
}

// UpdateDevices found in the network.
func (t *Tuner) UpdateDevices() error {
	ch := make(chan *OperationResult)

	for _, dev := range t.devices {
		go update(ch, dev)
	}

	errs := OperationErrors{}

	remaining := len(t.devices)

	for remaining != 0 {
		select {
		case result := <-ch:
			if result.finished {
				remaining--
			}

			if result.err != nil {
				errs = append(errs, NewOperationError(result.dev, result.err))
			}
		}
	}

	close(ch)

	return errs
}
