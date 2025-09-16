package shellygen1

import (
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
