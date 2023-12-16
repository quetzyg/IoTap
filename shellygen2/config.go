package shellygen2

import "github.com/Stowify/IoTune/device"

// The settings type increases flexibility when dealing with
// configuration payloads expected by the Shelly Gen2 device API.
type settings map[string]any

// Config implementation for the Shelly Gen2 driver.
type Config struct {
	Strategy *device.Strategy `json:"strategy,omitempty"`
	BLE      *settings        `json:"ble,omitempty"`
	Cloud    *settings        `json:"cloud,omitempty"`
	Eth      *settings        `json:"eth,omitempty"`
	Input    *[]*settings     `json:"input,omitempty"`
	MQTT     *settings        `json:"mqtt,omitempty"`
	Switch   *[]*settings     `json:"switch,omitempty"`
	Sys      *settings        `json:"sys,omitempty"`
	Wifi     *settings        `json:"wifi,omitempty"`
	WS       *settings        `json:"ws,omitempty"`
}

// Driver name of this Config implementation.
func (c *Config) Driver() string {
	return Driver
}

// Empty checks if the struct holding the configuration has a zero value.
func (c *Config) Empty() bool {
	return *c == Config{}
}
