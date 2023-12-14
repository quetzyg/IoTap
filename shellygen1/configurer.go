package shellygen1

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/Stowify/IoTune/device"
)

var paths = map[string]string{
	// Common HTTP API configuration endpoint paths
	// See: https://shelly-api-docs.shelly.cloud/gen1/#common-http-api
	"settings":         "settings",
	"settings_ap":      "settings/ap",
	"settings_sta":     "settings/sta",
	"settings_login":   "settings/login",
	"settings_cloud":   "settings/cloud",
	"settings_actions": "settings/actions",

	// Specific Shelly device API configuration endpoint paths
	// See: https://shelly-api-docs.shelly.cloud/gen1/#shelly1-shelly1pm
	"settings_relay":           "settings/relay/%d",
	"settings_power":           "settings/power/%d",
	"settings_ext_temperature": "settings/ext_temperature/%d",
	"settings_ext_humidity":    "settings/ext_humidity/%d",
	"settings_ext_switch":      "settings/ext_switch/%d",
}

// ConfigureRequests generates a slice of *http.Requests that are to be executed in order to configure an IoT device.
func (d *Device) ConfigureRequests(config device.Config) ([]*http.Request, error) {
	conf, match := config.(*Config)
	if !match {
		return nil, fmt.Errorf("%w: expected %q, got %q", device.ErrDriverMismatch, d.Driver(), config.Driver())
	}

	var requests []*http.Request

	confVal := reflect.Indirect(reflect.ValueOf(conf))

	for i := 0; i < confVal.Type().NumField(); i++ {
		setting := confVal.Field(i)

		// Skip nil setting pointers
		if setting.IsNil() {
			continue
		}

		// Current setting tag
		tag := strings.TrimSuffix(confVal.Type().Field(i).Tag.Get("json"), ",omitempty")
		path := paths[tag]

		switch params := setting.Interface().(type) {
		case *settings:
			r, err := request(d, path, params)
			if err != nil {
				return nil, err
			}

			requests = append(requests, r)

		case *[]*settings:
			for j, p := range *params {
				// Handle paths that require an index
				if strings.Contains(path, "%d") {
					path = fmt.Sprintf(path, j)
				}

				r, err := request(d, path, p)
				if err != nil {
					return nil, err
				}

				requests = append(requests, r)
			}
		}
	}

	// Ensure we reboot after applying all the settings
	if len(requests) > 0 {
		r, err := d.RebootRequest()
		if err != nil {
			return nil, err
		}

		requests = append(requests, r)
	}

	return requests, nil
}
