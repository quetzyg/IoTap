package shellygen2

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/quetzyg/IoTap/device"
)

// ConfigureRequests generates a slice of *http.Requests that are to be executed in order to configure an IoT device.
func (d *Device) ConfigureRequests(config device.Config) ([]*http.Request, error) {
	conf, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("%w: expected %q, got %q", device.ErrDriverMismatch, d.Driver(), config.Driver())
	}

	// Check if a configuration policy is set and enforce it
	if conf.Policy != nil && conf.Policy.IsExcluded(d) {
		return nil, device.ErrPolicyExcluded
	}

	var requests []*http.Request

	confVal := reflect.Indirect(reflect.ValueOf(conf))

	for i := range confVal.Type().NumField() {
		setting := confVal.Field(i)

		// Skip nil setting pointers
		if setting.IsNil() {
			continue
		}

		// Current setting tag
		tag := strings.TrimSuffix(confVal.Type().Field(i).Tag.Get("json"), ",omitempty")
		path := fmt.Sprintf("%s.SetConfig", tag)

		switch params := setting.Interface().(type) {
		case *settings:
			r, err := request(d, path, params)
			if err != nil {
				return nil, err
			}

			requests = append(requests, r)

		case *[]*settings:
			for _, p := range *params {
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
