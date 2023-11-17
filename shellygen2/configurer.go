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
		return nil, fmt.Errorf("%w: expected %q, got %q", device.ErrDriverMismatch, d.Driver(), cfg.Driver())
	}

	var requests []*http.Request

	if c.Settings != nil {
		if c.Settings.Input != nil {
			for _, in := range *c.Settings.Input {
				r, err := request(d, "Input.SetConfig", in)
				if err != nil {
					return nil, err
				}
				requests = append(requests, r)
			}
		}

		if c.Settings.Relay != nil {
			for _, rel := range *c.Settings.Relay {
				r, err := request(d, "Switch.SetConfig", rel)
				if err != nil {
					return nil, err
				}
				requests = append(requests, r)
			}
		}

		if c.Settings.Ethernet != nil {
			r, err := request(d, "Eth.SetConfig", c.Settings.Ethernet)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		if c.Settings.Wifi != nil {
			r, err := request(d, "Wifi.SetConfig", c.Settings.Wifi)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		if c.Settings.Bluetooth != nil {
			r, err := request(d, "BLE.SetConfig", c.Settings.Bluetooth)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		if c.Settings.Cloud != nil {
			r, err := request(d, "Cloud.SetConfig", c.Settings.Cloud)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		if c.Settings.MQTT != nil {
			r, err := request(d, "MQTT.SetConfig", c.Settings.MQTT)
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		// Reboot request
		r, err := request(d, "Shelly.Reboot", nil)
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}

	return requests, nil
}
