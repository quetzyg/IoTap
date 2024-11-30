package device

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Stowify/IoTune/httpclient"
)

const probeTimeout = time.Second * 8

// The Tuner type maintains a record of the Devices discovered during a network
// scan and has the capability to execute procedures on those devices.
type Tuner struct {
	probers []Prober
	config  Config
	script  *IoTScript
}

// SetConfig field value.
func (t *Tuner) SetConfig(cfg Config) {
	t.config = cfg
}

// SetScript field value.
func (t *Tuner) SetScript(scr *IoTScript) {
	t.script = scr
}

// NewTuner creates a new *Tuner instance.
func NewTuner(probers []Prober) *Tuner {
	return &Tuner{
		probers: probers,
	}
}

// Probe an IP address for a specific IoT device.
func Probe(client *http.Client, ip net.IP, prober Prober) (Resource, error) {
	r, dev, err := prober.Request(ip)
	if err != nil {
		return nil, err
	}

	err = httpclient.Dispatch(client, r, dev)

	var ue *url.Error
	if errors.As(err, &ue) {
		// Ignore timeouts, refused connections and other classic HTTP shenanigans,
		// since (NORMALLY!) it means there's no such device at the IP address.
		return nil, nil
	}

	if errors.Is(err, ErrUnexpected) {
		// Skip unexpected devices.
		return nil, nil
	}

	var je *json.SyntaxError
	if errors.As(err, &je) {
		// We found something, but it's not outputting valid JSON
		return nil, nil
	}

	return dev, err
}

// ProcedureResult encapsulates the outcome of a procedure executed on an IoT device.
// These can be related to various operations such as probing, updating, rebooting or configuring a device.
type ProcedureResult struct {
	dev Resource
	err error
}

// Failed checks if the ProcedureResult execution has failed.
func (pr *ProcedureResult) Failed() bool {
	return pr.err != nil
}

// Error interface implementation for ProcedureResult.
func (pr *ProcedureResult) Error() string {
	return fmt.Sprintf(
		"[%s] %s @ %s: %v\n",
		pr.dev.Driver(),
		pr.dev.ID(),
		pr.dev.IP(),
		pr.err,
	)
}

// probe an IP and send the result to a channel.
func probe(ch chan<- *ProcedureResult, ip net.IP, probers []Prober) {
	result := &ProcedureResult{}
	client := &http.Client{
		Timeout: probeTimeout,
	}

	for _, prober := range probers {
		dev, err := Probe(client, ip, prober)

		// Device found!
		if dev != nil {
			result.dev = dev
			break
		}

		if err != nil {
			result.err = &ProbeError{
				ip:  ip,
				err: err,
			}
		}
	}

	ch <- result
}

// Scan the network for IoT devices and return a Collection on success, error on failure.
func (t *Tuner) Scan(ips []net.IP) (Collection, error) {
	ch := make(chan *ProcedureResult)

	for _, ip := range ips {
		go probe(ch, ip, t.probers)
	}

	errs := Errors{}
	devices := Collection{}

	for range ips {
		result := <-ch
		if result.err != nil {
			errs = append(errs, result.err)
		}

		if result.dev != nil {
			devices = append(devices, result.dev)
		}
	}

	close(ch)

	if errs.Empty() {
		return devices, nil
	}

	return nil, errs
}

// procedure is a function type that encapsulates operations to be carried out on IoT devices.
type procedure func(tun *Tuner, res Resource, ch chan<- *ProcedureResult)

// Execute a procedure implementation on the found IoT devices.
func (t *Tuner) Execute(proc procedure, devices Collection) error {
	if len(devices) == 0 {
		return nil
	}

	ch := make(chan *ProcedureResult)

	for _, dev := range devices {
		go proc(t, dev, ch)
	}

	errs := Errors{}

	remaining := len(devices)
	for remaining != 0 {
		result := <-ch
		remaining--

		if result.Failed() {
			errs = append(errs, result)
		}
	}
	close(ch)

	if errs.Empty() {
		return nil
	}

	return errs
}
