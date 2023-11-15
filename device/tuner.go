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

// Tuner holds the Devices found during a network scan.
// It also has the ability to set their configuration.
type Tuner struct {
	probers []Prober
	config  Config
	devices Collection
}

// NewTuner creates a new tuner instance.
func NewTuner(probers []Prober, config Config) *Tuner {
	return &Tuner{
		probers: probers,
		config:  config,
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

// ProcedureResult encapsulates the outcome of a procedure executed on an IoT device.
// These can be related to various operations such as probing, updating, rebooting or configuring a device.
type ProcedureResult struct {
	dev Resource
	err error
}

// probe an IP and send the result to a channel.
func probe(ch chan<- *ProcedureResult, ip net.IP, probers []Prober) {
	result := &ProcedureResult{}

	for _, prober := range probers {
		dev, err := Probe(&http.Client{}, ip, prober)

		// Device found!
		if dev != nil {
			result.dev = dev
			break
		}

		if err != nil {
			result.err = NewProbeError(ip, err)
		}
	}

	ch <- result
	return
}

// The usable IP addresses of a /24 subnet.
const subnet24 = 254

// Scan the network with an IoT device prober.
func (t *Tuner) Scan(ip net.IP) error {
	// Cleanup before scanning
	t.devices = Collection{}

	ch := make(chan *ProcedureResult)

	for octet := byte(1); octet <= subnet24; octet++ {
		go probe(ch, net.IPv4(ip[0], ip[1], ip[2], octet), t.probers)
	}

	errs := Errors{}

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

// dispatch an HTTP request.
func dispatch(client *http.Client, r *http.Request) error {
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

// procedure is a function type designed to encapsulate operations to be carried out on an IoT device.
type procedure func(tun *Tuner, dev Resource, ch chan<- *ProcedureResult)

// Configure is a procedure implementation designed to apply configuration settings to an IoT device.
var Configure = func(tun *Tuner, dev Resource, ch chan<- *ProcedureResult) {
	rs, err := tun.config.MakeRequests(dev)
	if err != nil {
		ch <- &ProcedureResult{
			dev: dev,
			err: err,
		}
		return
	}

	client := &http.Client{}

	for _, r := range rs {
		if err = dispatch(client, r); err != nil {
			ch <- &ProcedureResult{
				dev: dev,
				err: err,
			}
			return
		}
	}

	ch <- &ProcedureResult{
		dev: dev,
	}
}

// Execute a procedure implementation on all IoT devices we have found.
func (t *Tuner) Execute(proc procedure) error {
	ch := make(chan *ProcedureResult)

	for _, dev := range t.devices {
		go proc(t, dev, ch)
	}

	errs := Errors{}

	remaining := len(t.devices)
	for remaining != 0 {
		select {
		case result := <-ch:
			remaining--

			if result.err != nil {
				errs = append(errs, NewOperationError(result.dev, result.err))
			}
		}
	}
	close(ch)

	return errs
}
