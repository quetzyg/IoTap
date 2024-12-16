package device

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
)

// The strategyMode determines whether the devices specified in the Strategy manifest should be configured or not.
// It supports two modes: blacklist (devices not to be configured) and whitelist (devices to be configured).
// This is helpful when dealing with multiple devices exhibiting minor configuration differences.
type strategyMode string

const (
	undefined strategyMode = ""
	blacklist strategyMode = "blacklist"
	whitelist strategyMode = "whitelist"
)

// Strategy to use when configuring IoT devices.
type Strategy struct {
	models  []string
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

	models, ok := m["models"].([]any)
	if ok {
		for _, str := range models {
			if model, ok := str.(string); ok {
				s.models = append(s.models, model)
			}
		}
	}

	// Currently, the MAC address unmarshalling logic has
	// to be done manually, due to Go's lack of support.
	// See: https://github.com/golang/go/issues/29678
	devices, ok := m["devices"].([]any)
	if ok {
		for _, dev := range devices {
			mac, err := net.ParseMAC(fmt.Sprint(dev))
			if err != nil {
				return err
			}

			s.devices = append(s.devices, mac)
		}
	}

	return nil
}

// Listed checks if a device model or MAC address is in the Strategy manifest.
func (s *Strategy) Listed(dev Resource) bool {
	for _, model := range s.models {
		if model == dev.Model() {
			return true
		}
	}

	for _, mac := range s.devices {
		if bytes.Equal(mac, dev.MAC()) {
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
