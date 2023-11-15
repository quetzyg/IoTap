package shellygen2

import (
	"fmt"
	"net/http"

	"github.com/Stowify/IoTune/device"
)

// ConfigureRequests generates a slice of *http.Requests that are to be executed in order to configure an IoT device.
func (d *Device) ConfigureRequests(cfg device.Config) ([]*http.Request, error) {
	c, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("driver mismatch, expected %s, got %s", d.Driver(), cfg.Driver())
	}

	var requests []*http.Request

	if c.Settings != nil {
		if c.Settings.Input != nil {
			for _, in := range *c.Settings.Input {
				r, err := makeRequest(d, "Input.SetConfig", in)
				if err != nil {
					return nil, err
				}
				requests = append(requests, r)
			}
		}

		if c.Settings.Relay != nil {
			for _, rel := range *c.Settings.Relay {
				r, err := makeRequest(d, "Switch.SetConfig", rel)
				if err != nil {
					return nil, err
				}
				requests = append(requests, r)
			}
		}

		if c.Settings.Ethernet != nil {
			r, err := makeRequest(d, "Eth.SetConfig", c.Settings.Ethernet)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		if c.Settings.Wifi != nil {
			r, err := makeRequest(d, "Wifi.SetConfig", c.Settings.Wifi)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		if c.Settings.Bluetooth != nil {
			r, err := makeRequest(d, "BLE.SetConfig", c.Settings.Bluetooth)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		if c.Settings.Cloud != nil {
			r, err := makeRequest(d, "Cloud.SetConfig", c.Settings.Cloud)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		if c.Settings.MQTT != nil {
			r, err := makeRequest(d, "MQTT.SetConfig", c.Settings.MQTT)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		// Reboot request
		r, err := makeRequest(d, "Shelly.Reboot", nil)
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}

	return requests, nil
}
