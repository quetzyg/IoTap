package iot

import (
	"net"
	"net/http"
)

type ProbeRequest func(net.IP) (*http.Request, Device, error)

// Device represents an IoT device.
type Device interface {
	Driver() string
	IP() net.IP
	ID() string
}

type Devices map[string]Device
