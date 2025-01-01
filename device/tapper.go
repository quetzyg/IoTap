package device

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/quetzyg/IoTap/httpclient"
)

const (
	probeTimeout  = time.Second * 8
	channelBuffer = 32
)

// Tapper knows how to tap into devices and execute tasks on them.
type Tapper struct {
	probers    []Prober
	config     Config
	cred       *Credentials
	auth       *AuthConfig
	deployment *Deployment
	transport  http.RoundTripper
}

// NewTapper creates a new *Tapper instance.
func NewTapper(probers []Prober) *Tapper {
	return &Tapper{
		probers: probers,
	}
}

// SetConfig implementation passed by the user.
func (t *Tapper) SetConfig(cfg Config) {
	t.config = cfg
}

// SetAuthConfig passed by the user.
func (t *Tapper) SetAuthConfig(auth *AuthConfig) {
	t.auth = auth
}

// SetDeployment passed by the user.
func (t *Tapper) SetDeployment(dep *Deployment) {
	t.deployment = dep
}

// probeIP for a specific IoT device.
func probeIP(prober Prober, client *http.Client, ip net.IP) (Resource, error) {
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

// probe an IP and return the probe result to a channel.
func (t *Tapper) probe(ch chan<- *ProcedureResult, client *http.Client, ip net.IP) {
	result := &ProcedureResult{}

	for _, prober := range t.probers {
		dev, err := probeIP(prober, client, ip)

		// Device found!
		if dev != nil {
			result.dev = dev

			if sec, ok := dev.(Securer); ok {
				sec.SetCredentials(t.cred)
			}
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
func (t *Tapper) Scan(ips []net.IP) (Collection, error) {
	ch := make(chan *ProcedureResult, channelBuffer)

	client := &http.Client{
		Transport: t.transport,
		Timeout:   probeTimeout,
	}

	for _, ip := range ips {
		go t.probe(ch, client, ip)
	}

	errs := Errors{}
	devices := Collection{}

	for range ips {
		result := <-ch
		if result.Failed() {
			errs = append(errs, result)
		}

		if result.dev != nil {
			devices = append(devices, result.dev)
		}
	}

	close(ch)

	if len(errs) == 0 {
		return devices, nil
	}

	return nil, errs
}

// Execute a procedure on a device collection.
func (t *Tapper) Execute(proc procedure, devices Collection) (int, error) {
	if devices.Empty() {
		return 0, nil
	}

	ch := make(chan *ProcedureResult, channelBuffer)

	for _, dev := range devices {
		go proc(t, dev, ch)
	}

	errs := Errors{}
	affected := 0

	for range devices {
		result := <-ch
		if !result.Failed() {
			affected++
			continue
		}

		// Skipped devices
		if errors.Is(result.err, ErrPolicyExcluded) {
			continue
		}

		errs = append(errs, result)
	}

	close(ch)

	if len(errs) == 0 {
		return affected, nil
	}

	return 0, errs
}
