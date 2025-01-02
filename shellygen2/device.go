package shellygen2

import (
	"encoding/json"
	"net"

	"github.com/quetzyg/IoTap/device"
)

const (
	// Driver name for this device implementation.
	Driver = "shelly_gen2"
)

// Device implementation for the Shelly Gen2 driver.
type Device struct {
	ip      net.IP
	mac     net.HardwareAddr
	name    string
	model   string
	secured bool

	Realm       string
	Generation  uint8
	Firmware    string
	Version     string
	VersionNext string
	AppName     string
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

// Secured returns true if the device requires authentication to be accessed, false otherwise.
func (d *Device) Secured() bool {
	return d.secured
}

// UnmarshalJSON implements the Unmarshaler interface.
func (d *Device) UnmarshalJSON(data []byte) error {
	// Unmarshal logic for the versioner implementation
	var fw struct {
		Src    string `json:"src"`
		Result struct {
			Stable struct {
				Version string `json:"version"`
			} `json:"stable"`
		} `json:"result"`
	}

	err := json.Unmarshal(data, &fw)
	if err != nil {
		return err
	}

	if fw.Src != "" {
		if fw.Result.Stable.Version != "" {
			d.VersionNext = fw.Result.Stable.Version
		}

		return nil
	}

	// Device unmarshal logic
	var dev struct {
		Name       *string `json:"name"`
		Realm      string  `json:"id"`
		MAC        string  `json:"mac"`
		Model      string  `json:"model"`
		Generation uint8   `json:"gen"`
		Firmware   string  `json:"fw_id"`
		Version    string  `json:"ver"`
		AppName    string  `json:"app"`
		Secured    bool    `json:"auth_en"`
	}

	err = json.Unmarshal(data, &dev)
	if err != nil {
		return err
	}

	// Different Shelly generations use different JSON field names,
	// but a Gen2 device should always have these fields populated.
	if dev.Model == "" && dev.Secured == false && dev.Firmware == "" {
		return device.ErrUnexpected
	}

	// Handle a potential nil name value
	d.name = "N/A"
	if dev.Name != nil {
		d.name = *dev.Name
	}

	d.Realm = dev.Realm

	d.mac, err = net.ParseMAC(device.Macify(dev.MAC))
	if err != nil {
		return err
	}

	d.model = dev.Model
	d.Generation = dev.Generation
	d.Firmware = dev.Firmware
	d.Version = dev.Version

	// Assume we're on the latest version, until we version the device.
	d.VersionNext = d.Version
	d.AppName = dev.AppName
	d.secured = dev.Secured

	return nil
}
