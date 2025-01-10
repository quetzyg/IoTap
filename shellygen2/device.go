package shellygen2

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/quetzyg/IoTap/device"
)

const (
	// Driver name for this device implementation.
	Driver = "shelly_gen2"

	// Vendor represents the name of the company that developed the device.
	Vendor = "Shelly"
)

// Device implementation for the Shelly Gen2 driver.
type Device struct {
	ip      net.IP
	mac     net.HardwareAddr
	name    string
	model   string
	secured bool
	cred    *device.Credentials

	Realm       string
	Gen         uint8
	Firmware    string
	Version     string
	VersionNext string
}

// IP address of the Device.
func (d *Device) IP() net.IP {
	return d.ip
}

// MAC address of the Device.
func (d *Device) MAC() net.HardwareAddr {
	return d.mac
}

// Name of the Device.
func (d *Device) Name() string {
	return d.name
}

// Model of the Device.
func (d *Device) Model() string {
	return d.model
}

// ID returns the Device's unique identifier.
func (d *Device) ID() string {
	return d.mac.String()
}

// Driver name of this Device implementation.
func (d *Device) Driver() string {
	return Driver
}

// Vendor represents the name of the company that developed the device.
func (d *Device) Vendor() string {
	return Vendor
}

// Generation represents the generation of this device.
func (d *Device) Generation() string {
	return fmt.Sprint(d.Gen)
}

// UnmarshalJSON implements the Unmarshaler interface.
func (d *Device) UnmarshalJSON(data []byte) error {
	// Unmarshal logic for the versioner implementation
	var fw struct {
		Result *struct {
			Stable struct {
				Version string `json:"version"`
			} `json:"stable"`
		} `json:"result"`
	}

	err := json.Unmarshal(data, &fw)
	if err != nil {
		return err
	}

	if fw.Result != nil {
		if fw.Result.Stable.Version != "" {
			d.VersionNext = fw.Result.Stable.Version
		}

		return nil
	}

	// Device unmarshal logic
	var dev struct {
		Name     *string `json:"name"`
		Realm    *string `json:"id"`
		MAC      *string `json:"mac"`
		Model    *string `json:"model"`
		Gen      *uint8  `json:"gen"`
		Firmware *string `json:"fw_id"`
		Version  *string `json:"ver"`
		Secured  *bool   `json:"auth_en"`
	}

	err = json.Unmarshal(data, &dev)
	if err != nil {
		return err
	}

	// Different Shelly generations use different JSON field names,
	// but a Gen2 device should always have these fields populated.
	if dev.Model == nil || dev.Secured == nil || dev.Firmware == nil || dev.Realm == nil || dev.Version == nil {
		return device.ErrUnexpected
	}

	// Default device name if unspecified
	d.name = "N/A"
	if dev.Name != nil {
		d.name = *dev.Name
	}

	d.Realm = *dev.Realm

	d.mac, err = net.ParseMAC(device.Macify(*dev.MAC))
	if err != nil {
		return err
	}

	d.model = *dev.Model
	d.Gen = *dev.Gen
	d.Firmware = *dev.Firmware
	d.Version = *dev.Version

	// Assume we're on the latest version, until we version the device.
	d.VersionNext = d.Version
	d.secured = *dev.Secured

	return nil
}
