package shellygen1

import (
	"encoding/json"
	"net"

	"github.com/quetzyg/IoTap/device"
)

const (
	// Driver name for this device implementation.
	Driver = "shellygen1"

	// Vendor represents the name of the company that developed the device.
	Vendor = "Shelly"

	// Generation of this device implementation.
	Generation = "1"
)

// Device implementation for the Shelly Gen1 driver.
type Device struct {
	cred         *device.Credentials
	name         string
	model        string
	Firmware     string
	FirmwareNext string
	ip           net.IP
	mac          net.HardwareAddr
	relays       uint8
	secured      bool
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

// Secured returns true if the device requires authentication to be accessed, false otherwise.
func (d *Device) Secured() bool {
	return d.secured
}

// versionUnmarshal checks if a firmware update is available for the device
// by extracting the "new_version" field from the JSON payload. If a valid
// firmware version is found, it updates the device's FirmwareNext field.
func (d *Device) versionUnmarshal(data []byte) error {
	var v struct {
		NewVersion *string `json:"new_version"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v.NewVersion != nil && *v.NewVersion != "" {
		d.FirmwareNext = *v.NewVersion
		return nil
	}

	return device.ErrUnexpected
}

// enrichUnmarshal attempts to update the device with additional metadata.
func (d *Device) enrichUnmarshal(data []byte) error {
	var v struct {
		Device struct {
			MAC *string `json:"mac"`
		} `json:"device"`
		Name *string `json:"name"`
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v.Device.MAC == nil {
		return device.ErrUnexpected
	}

	if v.Name != nil {
		d.name = *v.Name
	}

	return nil
}

// probeUnmarshal attempts to extract and parse fundamental device information
// from the JSON payload. This typically includes the device model, MAC address,
// security status, and firmware version. If the required fields are missing,
// it returns an error indicating an unexpected format.
func (d *Device) probeUnmarshal(data []byte) error {
	var v struct {
		Model    *string `json:"type"`
		MAC      *string `json:"mac"`
		Secured  *bool   `json:"auth"`
		Firmware *string `json:"fw"`
	}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	// Different Shelly generations use different JSON field names,
	// but a Gen1 device should always have these fields populated.
	if v.Model == nil || v.MAC == nil || /* v.Relays == nil || */ v.Secured == nil || v.Firmware == nil {
		return device.ErrUnexpected
	}

	d.mac, err = net.ParseMAC(device.Macify(*v.MAC))
	if err != nil {
		return err
	}

	d.model = *v.Model
	//	d.relays = *v.Relays

	// The /shelly endpoint for Gen1 devices does not provide
	// a name field, so we default to "N/A" and enrich later
	d.name = "N/A"

	// Assume we're on the latest version, until the device is versioned.
	d.FirmwareNext = *v.Firmware
	d.Firmware = *v.Firmware

	d.secured = *v.Secured

	return nil
}

// UnmarshalJSON implements the Unmarshaler interface.
func (d *Device) UnmarshalJSON(data []byte) error {
	if err := d.versionUnmarshal(data); err == nil {
		return nil
	}

	if err := d.enrichUnmarshal(data); err == nil {
		return nil
	}

	return d.probeUnmarshal(data)
}
