package shellygen2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Stowify/IoTune/internal/iot"
)

const (
	rpcPath = "rpc"
)

// Config implementation for the Shelly Gen2 driver.
type Config struct {
	Settings *Settings `json:"settings"`
}

// Settings for Shelly Gen2 devices.
type Settings struct {
	Input     *[]*input  `json:"input,omitempty"`
	Relay     *[]*relay  `json:"relay,omitempty"`
	Wifi      *wifi      `json:"wifi,omitempty"`
	Ethernet  *ethernet  `json:"ethernet,omitempty"`
	Bluetooth *bluetooth `json:"bluetooth,omitempty"`
	Cloud     *cloud     `json:"cloud,omitempty"`
	MQTT      *mqtt      `json:"mqtt,omitempty"`
}

// relay settings for Shelly Gen2 devices.
// Read more at: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Switch#configuration
// Some fields can have explicit null values, so they're not omitted
type relay struct {
	ID     int `json:"id"` // ID of the Switch component instance
	Config struct {
		Name                     *string `json:"name"`                                 // Name of the switch instance
		InMode                   *string `json:"in_mode,omitempty"`                    // Mode of the associated input: momentary, follow, flip, detached
		InitialState             *string `json:"initial_state,omitempty"`              // Output state to set on power_on: off, on, restore_last, match_input
		AutoOn                   *bool   `json:"auto_on,omitempty"`                    // True if the "Automatic ON" function is enabled, false otherwise
		AutoOnDelay              *int    `json:"auto_on_delay,omitempty"`              // Seconds to pass until the component is switched back on
		AutoOff                  *bool   `json:"auto_off,omitempty"`                   // True if the "Automatic OFF" function is enabled, false otherwise
		AutoOffDelay             *int    `json:"auto_off_delay,omitempty"`             // Seconds to pass until the component is switched back off
		AutoRecoverVoltageErrors *bool   `json:"autorecover_voltage_errors,omitempty"` // True if switch output state should be restored after over/under voltage error is cleared, false otherwise (shown if applicable)
		InputID                  *int    `json:"input_id,omitempty"`                   // ID of the Input component which controls the Switch. Applicable only to Pro1 and Pro1PM devices. Valid values: 0, 1
		PowerLimit               *int    `json:"power_limit,omitempty"`                // Limit (in Watts) over which overpower condition occurs (shown if applicable)
		VoltageLimit             *int    `json:"voltage_limit,omitempty"`              // Limit (in Volts) over which over voltage condition occurs (shown if applicable)
		UnderVoltageLimit        *int    `json:"undervoltage_limit,omitempty"`         // Limit (in Volts) under which under voltage condition occurs (shown if applicable)
		CurrentLimit             *int    `json:"current_limit,omitempty"`              // Number, limit (in Amperes) over which over current condition occurs (shown if applicable)
	} `json:"config"`
}

// input settings for Shelly Gen2 devices.
// https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Input#configuration
// Some fields can have explicit null values, so they're not omitted
type input struct {
	ID     int `json:"id"` // ID of the Input component instance
	Config struct {
		Name            *string `json:"name"`                    // Name of the input instance
		Type            *string `json:"type,omitempty"`          // Type of associated input: switch, button, analog (only if applicable)
		Invert          *bool   `json:"invert,omitempty"`        // True if the logical state of the associated input is inverted, false otherwise (switch and button type only)
		FactoryReset    *bool   `json:"factory_reset,omitempty"` // True if input-triggered factory reset option is enabled, false otherwise (switch and button type only)
		ReportThreshold *int    `json:"report_thr,omitempty"`    // Analog input report threshold in percent (analog type only)
	} `json:"config"`
}

// ethernet settings for Shelly Gen2 devices.
// Read more at: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Eth#configuration
// Some fields can have explicit null values, so they're not omitted
type ethernet struct {
	Config struct {
		Enable     *bool   `json:"enable,omitempty"`   // True if the configuration is enabled, false otherwise
		Ipv4Mode   *string `json:"ipv4mode,omitempty"` // IPv4 mode: dhcp, static
		IP         *string `json:"ip"`                 // IP to use when IPv4 mode is static
		Netmask    *string `json:"netmask"`            // Netmask to use when IPv4 mode is static
		Gateway    *string `json:"gw"`                 // Gateway to use when IPv4 mode is static
		Nameserver *string `json:"nameserver"`         // Nameserver to use when IPv4 mode is static
	} `json:"config"`
}

// wifi settings for Shelly Gen2 devices.
// Read more at: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/WiFi/#configuration
// Some fields can have explicit null values, so they're not omitted
type wifi struct {
	Config struct {
		AP       *wifiAP      `json:"ap,omitempty"`
		Station  *wifiStation `json:"sta,omitempty"`
		Station1 *wifiStation `json:"sta1,omitempty"`
		Roam     *wifiRoam    `json:"roam,omitempty"`
	} `json:"config"`
}

type wifiAP struct {
	Password      *string            `json:"pass"`                     // Password for the SSID
	IsOpen        *bool              `json:"is_open,omitempty"`        // True if the access point is open, false otherwise
	Enable        *bool              `json:"enable,omitempty"`         // True if the access point is enabled, false otherwise
	RangeExtender *wifiRangeExtender `json:"range_extender,omitempty"` // True if range extender functionality is enabled, false otherwise
}

type wifiRangeExtender struct {
	Enable *bool `json:"enable,omitempty"` // True if range extender functionality is enabled, false otherwise
}

type wifiStation struct {
	SSID       *string `json:"ssid,omitempty"`     // SSID of the network
	Password   *string `json:"pass"`               // Password for the SSID
	Enable     *bool   `json:"enable,omitempty"`   // True if the configuration is enabled, false otherwise
	Ipv4Mode   *string `json:"ipv4mode,omitempty"` // IPv4 mode: dhcp, static
	IP         *string `json:"ip"`                 // IP to use when IPv4 mode is static
	Netmask    *string `json:"netmask"`            // Netmask to use when IPv4 mode is static
	Gateway    *string `json:"gw"`                 // Gateway to use when IPv4 mode is static
	Nameserver *string `json:"nameserver"`         // Nameserver to use when IPv4 mode is static
}

type wifiRoam struct {
	RSSI     *int `json:"rssi_thr,omitempty"` // RSSI threshold - when reached will trigger the access point roaming. Default value: -80
	Interval *int `json:"interval,omitempty"` // Interval at which to scan for better access points. Enabled if set to positive number, disabled if set to 0. Default value: 60
}

// bluetooth settings for Shelly Gen2 devices.
// Read more at: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/BLE#configuration
type bluetooth struct {
	Config struct {
		Enable   *bool              `json:"enable,omitempty"` // True if bluetooth is enabled, false otherwise
		RPC      *bluetoothRPC      `json:"rpc,omitempty"`
		Observer *bluetoothObserver `json:"observer,omitempty"`
	} `json:"config"`
}

type bluetoothRPC struct {
	Enable *bool `json:"enable,omitempty"` // True if RPC service is enabled, false otherwise
}

type bluetoothObserver struct {
	Enable *bool `json:"enable,omitempty"` // True if BT LE observer is enabled, false otherwise
}

// cloud settings for Shelly Gen2 devices.
// Read more at: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Cloud/#configuration
type cloud struct {
	Config struct {
		Enable *bool `json:"enable,omitempty"` // True if cloud connection is enabled, false otherwise
	} `json:"config"`
}

// mqtt settings for Shelly Gen2 devices.
// Read more at: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Mqtt#configuration
// Some fields can have explicit null values, so they're not omitted
type mqtt struct {
	Config struct {
		Enable                *bool   `json:"enable,omitempty"`          // True if MQTT connection is enabled, false otherwise
		Server                *string `json:"server"`                    // Host name of the MQTT server. Can be followed by port number - host:port
		ClientID              *string `json:"client_id,omitempty"`       // Identifies each MQTT client that connects to an MQTT broker
		Username              *string `json:"user"`                      // MQTT username
		Password              *string `json:"pass"`                      // MQTT password
		SSLCA                 *string `json:"ssl_ca"`                    // Type of the TCP sockets
		TopicPrefix           *string `json:"topic_prefix,omitempty"`    // Prefix of the topics on which device publish/subscribe
		EnableRPC             *bool   `json:"enable_rpc,omitempty"`      // Enable RPC over MQTT
		RPCNotification       *bool   `json:"rpc_ntf,omitempty"`         // Enable RPC notifications (NotifyStatus and NotifyEvent) to be published
		StatusNotification    *bool   `json:"status_ntf,omitempty"`      // Enable publishing the complete component status
		UseClientCertificates *bool   `json:"use_client_cert,omitempty"` // Enable or disable usage of client certificates to use MQTT with encryption
		EnableControl         *bool   `json:"enable_control,omitempty"`  // Enable the MQTT control feature
	} `json:"config"`
}

// rpcRequest for the RPC endpoint.
type rpcRequest struct {
	ID         int    `json:"id"`
	Source     string `json:"src"`
	Method     string `json:"method"`
	Parameters any    `json:"params,omitempty"`
}

// Driver name of this Config implementation.
func (c *Config) Driver() string {
	return Driver
}

// MakeRequests returns a Shelly Gen2 specific HTTP request collection.
func (c *Config) MakeRequests(dev iot.Device) ([]*http.Request, error) {
	if dev.Driver() != c.Driver() {
		return nil, fmt.Errorf("device mismatch, expected %s, got %s", c.Driver(), dev.Driver())
	}

	var requests []*http.Request

	if c.Settings != nil {
		if c.Settings.Input != nil {
			for _, in := range *c.Settings.Input {
				r, err := makeRequest(dev, "Input.SetConfig", in)
				if err != nil {
					return nil, err
				}
				requests = append(requests, r)
			}
		}

		if c.Settings.Relay != nil {
			for _, rel := range *c.Settings.Relay {
				r, err := makeRequest(dev, "Switch.SetConfig", rel)
				if err != nil {
					return nil, err
				}
				requests = append(requests, r)
			}
		}

		if c.Settings.Ethernet != nil {
			r, err := makeRequest(dev, "Eth.SetConfig", c.Settings.Ethernet)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		if c.Settings.Wifi != nil {
			r, err := makeRequest(dev, "Wifi.SetConfig", c.Settings.Wifi)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		if c.Settings.Bluetooth != nil {
			r, err := makeRequest(dev, "BLE.SetConfig", c.Settings.Bluetooth)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		if c.Settings.Cloud != nil {
			r, err := makeRequest(dev, "Cloud.SetConfig", c.Settings.Cloud)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		if c.Settings.MQTT != nil {
			r, err := makeRequest(dev, "MQTT.SetConfig", c.Settings.MQTT)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		// Reboot request
		r, err := makeRequest(dev, "Shelly.Reboot", nil)
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}

	return requests, nil
}

// Empty checks if the struct holding the configuration has a zero value.
func (c *Config) Empty() bool {
	return *c == Config{}
}

// makeRequest for a Shelly Gen2 endpoint.
func makeRequest(dev iot.Device, method string, params any) (*http.Request, error) {
	req := &rpcRequest{
		Source: "IoTune",
		Method: method,
	}

	if params != nil {
		req.Parameters = params
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, buildURL(dev.IP(), rpcPath), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	r.Header.Set(iot.ContentTypeHeader, iot.JSONMimeType)

	return r, nil
}
