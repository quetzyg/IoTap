package device

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
)

// PolicyMode to apply when deploying to/configuring IoT devices.
// Two modes are supported: blacklist (devices to be excluded) and whitelist (devices to be included).
type PolicyMode int

const (
	PolicyModeUndefined PolicyMode = iota
	PolicyModeBlacklist
	PolicyModeWhitelist
)

var policyTypes = map[string]PolicyMode{
	"":          PolicyModeUndefined,
	"blacklist": PolicyModeBlacklist,
	"whitelist": PolicyModeWhitelist,
}

// UnmarshalJSON implements the Unmarshaler interface.
func (pm *PolicyMode) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	var ok bool
	if *pm, ok = policyTypes[s]; !ok {
		return fmt.Errorf("%w: %s", errPolicyModeInvalid, s)
	}

	if *pm == PolicyModeUndefined {
		return errPolicyModeUndefined
	}

	return nil
}

// Policy to apply when deploying to/configuring IoT devices.
type Policy struct {
	Mode    PolicyMode         `json:"mode"`
	Names   []string           `json:"names"`
	Models  []string           `json:"models"`
	Devices []net.HardwareAddr `json:"devices"`
}

// Contains checks whether a device model or MAC address exists in the Policy.
func (p *Policy) Contains(dev Resource) bool {
	for _, name := range p.Names {
		if regexp.MustCompile(name).MatchString(dev.Name()) {
			return true
		}
	}

	for _, model := range p.Models {
		if regexp.MustCompile(model).MatchString(dev.Model()) {
			return true
		}
	}

	for _, mac := range p.Devices {
		if bytes.Equal(mac, dev.MAC()) {
			return true
		}
	}

	return false
}

// IsExcluded determines whether a device should be excluded based on the Policy mode
// (blacklist or whitelist) and whether the device is contained within the Policy.
func (p *Policy) IsExcluded(dev Resource) bool {
	switch p.Mode {
	case PolicyModeBlacklist:
		return p.Contains(dev)

	case PolicyModeWhitelist:
		return !p.Contains(dev)

	default:
		panic("invalid policy mode")
	}
}

// UnmarshalJSON implements the Unmarshaler interface.
func (p *Policy) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Type    PolicyMode `json:"mode"`
		Names   []string   `json:"names"`
		Models  []string   `json:"models"`
		Devices []string   `json:"devices"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if tmp.Type == PolicyModeUndefined {
		return errPolicyModeUndefined
	}

	p.Mode = tmp.Type
	p.Names = tmp.Names
	p.Models = tmp.Models

	// Currently, the MAC address unmarshalling logic has
	// to be done manually, due to Go's lack of support.
	// See: https://github.com/golang/go/issues/29678
	for _, dev := range tmp.Devices {
		mac, err := net.ParseMAC(dev)
		if err != nil {
			return err
		}

		p.Devices = append(p.Devices, mac)
	}

	return nil
}
