package shellygen1

import (
	"encoding/json"
	"net"

	"github.com/quetzyg/IoTap/device"
)

const (
	// Driver name for this device implementation.
	Driver = "shelly_gen1"

	// Vendor represents the name of the company that developed the device.
	Vendor = "Shelly"

	// Generation of this device implementation.
	Generation = "1"
)

// Device implementation for the Shelly Gen1 driver.
type Device struct {
	ip      net.IP
	mac     net.HardwareAddr
	name    string
	model   string
	secured bool
	cred    *device.Credentials

	Firmware     string
	FirmwareNext string
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
	return Generation
}

// UnmarshalJSON implements the Unmarshaler interface.
func (d *Device) UnmarshalJSON(data []byte) error {
	// Unmarshal logic for the versioner implementation
	var fw struct {
		NewVersion *string `json:"new_version"`
	}

	err := json.Unmarshal(data, &fw)
	if err != nil {
		return err
	}

	if fw.NewVersion != nil {
		if *fw.NewVersion != "" {
			d.FirmwareNext = *fw.NewVersion
		}

		return nil
	}

	// Device unmarshal logic
	var dev struct {
		Model    *string `json:"type"`
		MAC      *string `json:"mac"`
		Secured  *bool   `json:"auth"`
		Firmware *string `json:"fw"`
	}

	err = json.Unmarshal(data, &dev)
	if err != nil {
		return err
	}

	// Different Shelly generations use different JSON field names,
	// but a Gen1 device should always have these fields populated.
	if dev.Model == nil || dev.Secured == nil || dev.Firmware == nil {
		return device.ErrUnexpected
	}

	d.model = *dev.Model

	d.mac, err = net.ParseMAC(device.Macify(*dev.MAC))
	if err != nil {
		return err
	}

	// Unfortunately, the /shelly endpoint for Gen1 devices does not provide
	// a field for the device name, so we currently populate it with "N/A"
	d.name = "N/A"
	d.Firmware = *dev.Firmware

	// Assume we're on the latest version, until we version the device.
	d.FirmwareNext = d.Firmware

	d.secured = *dev.Secured

	return nil
}
