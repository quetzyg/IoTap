package device

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	iotune "github.com/Stowify/IoTune"
)

const probeTimeout = time.Second * 8

// The Tuner type maintains a record of the Devices discovered during a network
// scan and has the capability to execute procedures on those devices.
type Tuner struct {
	probers []Prober
	config  Config
	script  *IoTScript
	devices Collection
}

// SetConfig field value.
func (t *Tuner) SetConfig(cfg Config) {
	t.config = cfg
}

// SetScript field value.
func (t *Tuner) SetScript(scr *IoTScript) {
	t.script = scr
}

// NewTuner creates a new tuner instance.
func NewTuner(probers []Prober) *Tuner {
	return &Tuner{
		probers: probers,
		devices: Collection{},
	}
}

// Probe an IP address for a specific IoT device.
func Probe(client *http.Client, ip net.IP, prober Prober) (Resource, error) {
	r, dev, err := prober.ProbeRequest(ip)
	if err != nil {
		return nil, err
	}

	err = iotune.Dispatch(client, r, dev)

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
			result.err = NewProbeError(ip, err)
		}
	}

	ch <- result
}

// Scan the network with an IoT device prober.
func (t *Tuner) Scan(ips []net.IP) error {
	// Cleanup before scanning
	t.devices = Collection{}

	ch := make(chan *ProcedureResult)

	for _, ip := range ips {
		go probe(ch, ip, t.probers)
	}

	errs := Errors{}

	for i := 0; i < len(ips); i++ {
		result := <-ch
		if result.err != nil {
			errs = append(errs, result.err)
		}

		if result.dev != nil {
			t.devices[result.dev.ID()] = result.dev
		}
	}

	close(ch)

	return errs
}

// Devices that were found during the network scan.
func (t *Tuner) Devices() Collection {
	return t.devices
}

// procedure is a function type designed to encapsulate operations to be carried out on an IoT device.
type procedure func(tun *Tuner, dev Resource, ch chan<- *ProcedureResult)

// Execute a procedure implementation on all IoT devices we have found.
func (t *Tuner) Execute(proc procedure) error {
	ch := make(chan *ProcedureResult)

	for _, dev := range t.devices {
		go proc(t, dev, ch)
	}

	errs := Errors{}

	remaining := len(t.devices)
	for remaining != 0 {
		result := <-ch
		remaining--

		if result.err != nil {
			errs = append(errs, NewOperationError(result.dev, result.err))
		}
	}
	close(ch)

	return errs
}
