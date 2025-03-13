package shellygen2

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/quetzyg/IoTap/device"
)

const (
	// Driver name for this device implementation.
	Driver = "shellygen2"

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

// Secured returns true if the device requires authentication to be accessed, false otherwise.
func (d *Device) Secured() bool {
	return d.secured
}

// versionUnmarshal checks if a firmware update is available for the device
// by extracting the "new_version" field from the JSON payload. If a valid
// firmware version is found, it updates the device's FirmwareNext field.
func (d *Device) versionUnmarshal(data []byte) error {
	var v struct {
		Result *struct {
			Stable struct {
				Version string `json:"version"`
			} `json:"stable"`
		} `json:"result"`
	}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	if v.Result != nil {
		if v.Result.Stable.Version != "" {
			d.VersionNext = v.Result.Stable.Version
		}

		// We ignore other versions (i.e. beta)

		return nil
	}

	return device.ErrUnexpected
}

// probeUnmarshal attempts to extract and parse fundamental device information
// from the JSON payload. This typically includes the device model, MAC address,
// security status, and firmware version. If the required fields are missing,
// it returns an error indicating an unexpected format.
func (d *Device) probeUnmarshal(data []byte) error {
	var v struct {
		Name     *string `json:"name"`
		Realm    *string `json:"id"`
		MAC      *string `json:"mac"`
		Model    *string `json:"model"`
		Gen      *uint8  `json:"gen"`
		Firmware *string `json:"fw_id"`
		Version  *string `json:"ver"`
		Secured  *bool   `json:"auth_en"`
	}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	// Different Shelly generations use different JSON field names,
	// but a Gen2 device should always have these fields populated.
	if v.Realm == nil || v.MAC == nil || v.Model == nil ||
		v.Gen == nil || v.Firmware == nil || v.Version == nil ||
		v.Secured == nil {
		return device.ErrUnexpected
	}

	d.mac, err = net.ParseMAC(device.Macify(*v.MAC))
	if err != nil {
		return err
	}

	// Default device name if unspecified
	d.name = "N/A"
	if v.Name != nil {
		d.name = *v.Name
	}

	d.Realm = *v.Realm
	d.model = *v.Model
	d.Gen = *v.Gen
	d.Firmware = *v.Firmware

	// Assume we're on the latest version, until the device is versioned.
	d.VersionNext = *v.Version
	d.Version = *v.Version

	d.secured = *v.Secured

	return nil
}

// UnmarshalJSON implements the Unmarshaler interface.
func (d *Device) UnmarshalJSON(data []byte) error {
	if err := d.versionUnmarshal(data); err == nil {
		return nil
	}

	return d.probeUnmarshal(data)
}
