package device

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	iotune "github.com/Stowify/IoTune"
)

// Configurer is an interface that provides a standard way to configure IoT devices.
type Configurer interface {
	ConfigureRequests(Config) ([]*http.Request, error)
}

// The strategyMode determines whether the devices specified in the Strategy manifest should be configured or not.
// It supports two modes: blacklist (devices not to be configured) and whitelist (devices to be configured).
// This is helpful when dealing with multiple devices exhibiting minor configuration differences.
type strategyMode string

const (
	undefined strategyMode = ""
	blacklist              = "blacklist"
	whitelist              = "whitelist"
)

// Strategy to use when configuring IoT devices.
type Strategy struct {
	mode    strategyMode
	devices []net.HardwareAddr
}

// UnmarshalJSON implements the Unmarshaler interface.
func (s *Strategy) UnmarshalJSON(data []byte) error {
	var m map[string]any
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	mode, ok := m["mode"].(string)
	if !ok {
		return errStrategyModeUndefined
	}

	switch strategyMode(mode) {
	case blacklist, whitelist:
		s.mode = strategyMode(mode)
	case undefined:
		return errStrategyModeUndefined
	default:
		return fmt.Errorf("%w: %s", errStrategyModeInvalid, s.mode)
	}

	// Currently, the MAC address unmarshalling logic has
	// to be done manually, due to Go's lack of support.
	// See: https://github.com/golang/go/issues/29678
	devices, ok := m["devices"].([]any)
	if !ok {
		return nil
	}

	for _, dev := range devices {
		mac, err := net.ParseMAC(Macify(fmt.Sprint(dev)))
		if err != nil {
			return err
		}

		s.devices = append(s.devices, mac)
	}

	return nil
}

// Listed checks if a device MAC address is in the Strategy manifest.
func (s *Strategy) Listed(dev Resource) bool {
	for _, mac := range s.devices {
		if mac.String() == dev.MAC().String() {
			return true
		}
	}

	return false
}

// Excluded checks if a device is to be excluded from being configured.
func (s *Strategy) Excluded(dev Resource) bool {
	if s.mode == blacklist {
		return s.Listed(dev)
	}

	return !s.Listed(dev)
}

// Configure is a procedure implementation designed to apply configuration settings to an IoT device.
var Configure = func(tun *Tuner, dev Resource, ch chan<- *ProcedureResult) {
	res, ok := dev.(Configurer)
	if !ok {
		ch <- &ProcedureResult{
			err: fmt.Errorf("%w: configure", ErrUnsupportedProcedure),
			dev: dev,
		}
		return
	}

	rs, err := res.ConfigureRequests(tun.config)
	if err != nil {
		ch <- &ProcedureResult{
			dev: dev,
			err: err,
		}
		return
	}

	// A device without requests means it was
	// excluded as part of the strategy logic
	if len(rs) == 0 {
		ch <- &ProcedureResult{
			dev: dev,
			err: errStrategyExcludedDevice,
		}
		return
	}

	client := &http.Client{}

	for _, r := range rs {
		if err = iotune.Dispatch(client, r, nil); err != nil {
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
