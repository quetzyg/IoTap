package shellygen2

import (
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
	cred        *device.Credentials
	name        string
	model       string
	Realm       string
	Firmware    string
	Version     string
	VersionNext string
	ip          net.IP
	mac         net.HardwareAddr
	secured     bool
	Gen         uint8
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
